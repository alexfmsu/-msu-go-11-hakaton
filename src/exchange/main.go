package main

import (
	//"_exchange/brockerLogic"
	"exchange/GRPCExchangeServer"
	_ "exchange/brokerInterface"
	"exchange/marketLogic"
	_ "exchange/statistics"
	_ "exchange/structs"
	_ "exchange/transactionStorage"
	_ "exchange/transactionTypes"
	"fmt"
)

func main() {
	fmt.Println("Exchange started")

	marketLogic.OrdersChannelIn, marketLogic.OrdersChannelOut = marketLogic.InitCore()
	marketLogic.Init()
	// Initialize all here
	GRPCExchangeServer.InitServer()

	// Infinite loop
	// Don't remove
	// for {
	//
	// }
}

//PanicOnErr panics on error
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
