package structs

import (
	s "exchange/transactionStorage"
	types "exchange/transactionTypes"
)

type TickerStock struct {
	Sell *s.SortedOrders
	Buy  *s.SortedOrders
}

var Stock map[string]TickerStock
var TickerNames []string

func init() {
	Stock = map[string]TickerStock{
		"RIM7":         NewTickerStock(),
		"SPFB.Si-6.17": NewTickerStock(),
	}

	TickerNames = append(TickerNames, "RIM7")
	TickerNames = append(TickerNames, "SPFB.Si-6.17")
}

func NewTickerStock() TickerStock {
	t := TickerStock{}
	t.Buy = s.NewSortedOrders(types.HasMorePriorityForBuy)
	t.Sell = s.NewSortedOrders(types.HasMorePriorityForSell)
	return t
}
