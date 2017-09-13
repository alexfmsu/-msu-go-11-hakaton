package marketLogic

import (
	// server "exchange/GRPCExchangeServer"
	"exchange/structs"
	ts "exchange/transactionStorage"
	t "exchange/transactionTypes"
	"fmt"
	"time"
)

var (
	OrdersChannelIn  chan *t.Order
	OrdersChannelOut chan *t.Order
)

func InitCore() (chan *t.Order, chan *t.Order) {
	in := make(chan *t.Order)
	out := make(chan *t.Order)
	go coreFunction(in, out)
	return in, out
}

func coreFunction(in, out chan *t.Order) {
	for order := range in {
		var tickerStock structs.TickerStock
		var ok bool

		if order.Amount <= 0 || order.Price <= 0.0 {
			continue
		}

		order.Time = time.Now()

		if tickerStock, ok = structs.Stock[order.Ticker]; !ok {
			continue
		}

		if order.Action == t.ActionRemove {
			var storage *ts.SortedOrders
			switch order.Type {
			case t.Sell:
				storage = tickerStock.Sell
			case t.Buy:
				storage = tickerStock.Buy
			}
			err := storage.Delete(order.BrockerID, order.OrderID)
			if err != nil {
				fmt.Println("Not exist", t.StatusNotExist)
				order.Status = t.StatusNotExist
			} else {
				fmt.Println("Removed")
				order.Status = t.StatusRemoved
			}
			out <- order
			continue
		}

		if !order.Instant {
			fmt.Println("Stock transaction")
		}

		switch order.Type {
		case t.Sell:
			trySell(order, tickerStock, out)
		case t.Buy:
			tryBuy(order, tickerStock, out)
		}
	}
}

func trySell(orderForSell *t.Order, tickerStock structs.TickerStock, out chan *t.Order) {
	for {
		orderForBuy, err := tickerStock.Buy.Peek()
		if err != nil {
			break
		}
		if orderForSell.Price > orderForBuy.Price {
			break
		}

		minAmount := orderForBuy.Amount
		if orderForSell.Amount < minAmount {
			minAmount = orderForSell.Amount
		}

		orderForBuy.Amount -= minAmount
		orderForSell.Amount -= minAmount

		orderForBuy.Partial = orderForBuy.Amount != 0
		orderForSell.Partial = orderForSell.Amount != 0

		if !orderForSell.Instant {
			out <- orderForSell
		}
		out <- orderForBuy

		if orderForBuy.Amount == 0 {
			tickerStock.Buy.Pop()
		}

		if orderForSell.Amount == 0 {
			break
		}
	}

	if orderForSell.Amount > 0 && !orderForSell.Instant {
		tickerStock.Sell.Insert(orderForSell)
	}
}

func tryBuy(orderForBuy *t.Order, tickerStock structs.TickerStock, out chan *t.Order) {
	for {
		orderForSell, err := tickerStock.Sell.Peek()
		if err != nil {
			break
		}
		if orderForSell.Price > orderForBuy.Price {
			break
		}

		minAmount := orderForBuy.Amount
		if orderForSell.Amount < minAmount {
			minAmount = orderForSell.Amount
		}

		orderForBuy.Amount -= minAmount
		orderForSell.Amount -= minAmount

		orderForBuy.Partial = orderForBuy.Amount != 0
		orderForSell.Partial = orderForSell.Amount != 0

		if !orderForBuy.Instant {
			out <- orderForBuy
		}
		out <- orderForSell

		if orderForSell.Amount == 0 {
			tickerStock.Sell.Pop()
		}

		if orderForBuy.Amount == 0 {
			break
		}
	}

	if orderForBuy.Amount > 0 && !orderForBuy.Instant {
		tickerStock.Buy.Insert(orderForBuy)
	}
}
