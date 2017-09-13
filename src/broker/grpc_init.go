package main

import (
	cfg "broker/config"
	"context"
	"database/sql"
	bk "exchange/brokerInterface"

	"google.golang.org/grpc"
)

var grcpConn *grpc.ClientConn
var exchangeInterface bk.BrokerExcangeInterfaceClient
var requestStream bk.BrokerExcangeInterface_SendTransactionRequestStreamClient

func initGRPC() {
	var err error
	grcpConn, err = grpc.Dial(configData.Host+":"+configData.GrpcPort, grpc.WithInsecure())
	checkErr(err)
	exchangeInterface = bk.NewBrokerExcangeInterfaceClient(grcpConn)
	requestStream, err = exchangeInterface.SendTransactionRequestStream(context.Background())
	checkErr(err)
	for _, broker := range configData.Brokers[:1] {
		go func(broker cfg.Broker) {
			stat, err := exchangeInterface.GetStatistics(context.Background(), &bk.Broker{BrokerID: broker.ID})
			checkErr(err)
			for {
				price, err := stat.Recv()
				if err != nil {
					break
				}
				hist, ok := history[broker.ID]
				if !ok {
					hist.Data = map[string][]bk.OHLCV{}
				}
				for tickerName, data := range price.Ticker {
					tmp := hist.Data[tickerName]
					sliceIdx := 0
					if len(tmp) > configData.Timeout {
						sliceIdx = 1
					}
					tmp = append(tmp[sliceIdx:], *data)
					hist.Data[tickerName] = tmp
				}
				history[broker.ID] = hist
			}
		}(broker)
		go func(broker cfg.Broker) {
			resultGet, err := exchangeInterface.GetTransactionResultStream(context.Background(), &bk.Broker{BrokerID: broker.ID})
			checkErr(err)
			for {
				transAct, err := resultGet.Recv()
				checkErr(err)
				status := transAct.GetStatus()
				ordID := transAct.GetOrderID()
				if status != 2 && status != 4 {
					_, err = db.Exec(
						"UPDATE `positions` SET `status` = ? WHERE id = ?",
						transAct.GetStatus(),
						ordID,
					)
					checkErr(err)
				} else {
					var amount int
					err := db.QueryRow(
						"SELECT `amount` FROM `positions` WHERE `id`=?",
						ordID,
					).Scan(&amount)
					if err != sql.ErrNoRows {
						checkErr(err)
					}
					_, err = db.Exec(
						"INSERT INTO `orders_history` (`time`, `user_id`, `ticker`, `amount`, `price`, `type`)"+
							"VALUES (?, ?, ?, ?, ?, ?) ",
						transAct.GetTime(), transAct.GetClientID(), transAct.GetTicker(),
						amount, transAct.GetPrice(), transAct.GetTransactionType(),
					)
					checkErr(err)
					_, err = db.Exec(
						"DELETE FROM `positions` WHERE `id` = ?",
						ordID,
					)
					checkErr(err)
				}
			}
		}(broker)
	}
}
