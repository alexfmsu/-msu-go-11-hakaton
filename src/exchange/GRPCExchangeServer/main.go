package GRPCExchangeServer

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	bk "exchange/brokerInterface"
	logic "exchange/marketLogic"
	st "exchange/statistics"
	"exchange/structs"
	tt "exchange/transactionTypes"
	_ "golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// GRPCExchangeServerType - Server type
type GRPCExchangeServerType struct{}

var TransactionResultChannels map[int32]chan bk.TransactionRequestType

//////////////////////   Stream handlers  ///////////////////////////////

func SendResults(out chan *tt.Order) {
	for order := range out {
		if (order.Status != tt.StatusNotExist) && (order.Status != tt.StatusRemoved) {
			order.Status = tt.StatusCompleted
		}
		if order.Partial {
			order.Status = tt.StatusPartialCompleted
		}
		SendTransactionResult(*order)
	}
}

// Call this to send transaction result
func SendTransactionResult(order tt.Order) error {
	_, exist := TransactionResultChannels[int32(order.BrockerID)]
	if !exist {
		return errors.New("Broker not connected")
	}

	transaction := ConvertOrderToTransactionRequest(order)

	TransactionResultChannels[transaction.BrokerID] <- transaction
	return nil
}

// A client-to-server streaming RPC.
// New transaction requests from clients
// should be sent to exchange by broker through this channel
func (s *GRPCExchangeServerType) SendTransactionRequestStream(stream bk.BrokerExcangeInterface_SendTransactionRequestStreamServer) error {
	var request *bk.TransactionRequestType
	var err error
	for {
		request, err = stream.Recv()
		if err != nil {
			return err
		}

		// In case of errors status will be updated
		var order = ConvertTransactionRequestToOrder(*request)

		if order.Status < 0 {
			fmt.Println("Got transaction request from client", request.ClientID, "from", request.BrokerID, "with error", order.Status, " and action", request.Action)
		} else {
			fmt.Println("Got transaction request from client", request.ClientID, " type ", request.TransactionType, "with status ", order.Status, " and action", request.Action)
			logic.OrdersChannelIn <- &order
			order.Status = tt.StatusTransactionAdded
		}

		// Send request with updated status
		err := SendTransactionResult(order)

		if err != nil {
			fmt.Println(err.Error())
		}

	}
	return nil
}

// A server-to-client streaming RPC.
// Transaction results from the exchange
// should be sent back to clients through this channel
func (s *GRPCExchangeServerType) GetTransactionResultStream(bro *bk.Broker, stream bk.BrokerExcangeInterface_GetTransactionResultStreamServer) error {
	var request bk.TransactionRequestType
	//TransactionResultChannel = make(chan bk.TransactionRequestType)
	TransactionResultChannels[bro.BrokerID] = make(chan bk.TransactionRequestType)
	defer delete(TransactionResultChannels, bro.BrokerID)

	fmt.Println(bro.BrokerID, "Connected")

	for {
		// Don't modify
		// Use SendTransactionResult() to send result
		request = <-TransactionResultChannels[bro.BrokerID]
		// Send result
		stream.Send(&request)
	}

	return nil
}

// A server-to-client streaming RPC.
// Send statistics to brocker evry second
// To get statistics send BrockerID
func (s *GRPCExchangeServerType) GetStatistics(_ *bk.Broker, stream bk.BrokerExcangeInterface_GetStatisticsServer) error {
	var gStatMap bk.StatisticsType
	gStatMap.Ticker = make(map[string]*bk.OHLCV)
	var sStat st.OHLCV
	for {
		// Note: this finction will be called for all connected brokers every second
		// Put obtained statistics in sStat

		for _, ticker := range structs.TickerNames {
			var gStat bk.OHLCV
			sStat = st.GetStatistics(ticker)
			gStat = ConvertStOHLCVToGrpcOHLCV(sStat)
			gStatMap.Ticker[ticker] = &gStat
		}

		/*
			for _, ticker := range structs.TickerNames {
				fmt.Println(ticker, *gStatMap.Ticker[ticker])
			}
		*/

		stream.Send(&gStatMap)

		// Delay
		timer := time.NewTimer(time.Second * 1)
		<-timer.C
	}

	return nil
}

func InitServer() {
	TransactionResultChannels = make(map[int32]chan bk.TransactionRequestType)

	go SendResults(logic.OrdersChannelOut)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	bk.RegisterBrokerExcangeInterfaceServer(s, &GRPCExchangeServerType{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//////////////////////   Type conversion functions  ///////////////////////////////

func ConvertTransactionRequestToOrder(request bk.TransactionRequestType) tt.Order {
	var order tt.Order

	order.OrderID = request.OrderID
	switch request.TransactionType {
	case "sell":
		order.Type = tt.Sell
	case "buy":
		order.Type = tt.Buy
	default:
		order.Status = tt.StatusInvalidType
	}
	order.Price = float64(request.Price)
	order.Time = time.Unix(int64(request.Time), 0)
	order.Amount = int(request.Amount)
	order.BrockerID = int(request.BrokerID)
	order.ClientID = int(request.ClientID)
	order.Ticker = string(request.Ticker)
	order.Partial = bool(request.Partial)

	switch request.Action {
	case "add":
		order.Action = tt.ActionAdd
	case "remove":
		order.Action = tt.ActionRemove
	case "get":
		order.Action = tt.ActionGet
	default:
		order.Status = tt.StatusInvalidAction
	}

	// Check ticker
	cnt := 0
	for _, ticker := range structs.TickerNames {
		if ticker == order.Ticker {
			cnt++
		}
	}
	if cnt == 0 {
		order.Status = tt.StatusInvalidTicker
	}

	return order
}

func ConvertOrderToTransactionRequest(order tt.Order) bk.TransactionRequestType {
	var request bk.TransactionRequestType

	switch order.Type {
	case tt.Sell:
		request.TransactionType = "sell"
	case order.Type:
		request.TransactionType = "buy"
	}

	switch order.Action {
	case tt.ActionAdd:
		request.Action = "add"
	case tt.ActionRemove:
		request.Action = "remove"

	}
	request.OrderID = order.OrderID
	request.Price = float32(order.Price)
	request.Time = int32(order.Time.Unix())
	request.Amount = int32(order.Amount)
	request.BrokerID = int32(order.BrockerID)
	request.ClientID = int32(order.ClientID)
	request.Ticker = string(order.Ticker)
	request.Partial = bool(order.Partial)
	request.Status = order.Status

	return request
}

func ConvertStOHLCVToGrpcOHLCV(sStatistics st.OHLCV) bk.OHLCV {
	var gStatistics bk.OHLCV

	gStatistics.ID = int64(sStatistics.ID)
	gStatistics.Time = int32(sStatistics.Time.Unix())
	gStatistics.Interval = int32(sStatistics.Interval)
	gStatistics.Open = float32(sStatistics.Open)
	gStatistics.High = float32(sStatistics.High)
	gStatistics.Low = float32(sStatistics.Low)
	gStatistics.Close = float32(sStatistics.Close)
	gStatistics.Ticker = string(sStatistics.Ticker)
	return gStatistics
}
