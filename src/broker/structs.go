package main

import (
	cfg "broker/config"
	bk "exchange/brokerInterface"
)

// Deal ...
type Deal struct {
	ID     int     `json:"id"`
	Time   int     `json:"time"`
	UserID int     `json:"-"`
	Ticker string  `json:"ticker"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
	Status string  `json:"status"`
}

// User ...
type User struct {
	ID      int        `json:"id"`
	Name    string     `json:"name"`
	Balance float64    `json:"balance"`
	Token   string     `json:"token"`
	Broker  cfg.Broker `json:"broker"`
}

// History ...
type History struct {
	BrokerID int32
	Data     map[string][]bk.OHLCV
}

// Response is general response type
type Response struct {
	Data  interface{} `json:"data"`
	UserS *User       `json:"user"`
	Error string      `json:"error"`
}
