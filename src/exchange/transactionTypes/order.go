package transactionTypes

import (
	"time"
)

//OrderType - represents order type: Sell or Buy
type OrderType int

//Action - что должны сделать с order, добавить или удалить
type Action int

const (
	//Sell - sell order
	Sell OrderType = iota
	//Buy - buy order
	Buy
)

const (
	ActionAdd    = 0 // Add - add order
	ActionRemove = 1 // Delete - delete order
	ActionGet    = 2 // Get - Get transaction from storage without removing (for information updating)
)

const (
	StatusOK               = 0  // статус не определён (ставится клиентом при отправке запроса)
	StatusTransactionAdded = 1  // запрос на транзакцию добавлен
	StatusCompleted        = 2  // транзакция совершена полностью
	StatusPartialCompleted = 3  // транзакция совершена частично
	StatusRemoved          = 4  // запрос на транзакцию удалён
	StatusInvalidTicker    = -1 // ошибка, несуществующий Ticker.
	StatusInvalidType      = -2 // ошибка, несуществующий Ticker.
	StatusInvalidAction    = -3 // ошибка, несуществующий Ticker.
	StatusNotExist         = -4 // ошибка, записи не существует
)

//Order -
type Order struct {
	Type      OrderType
	Price     float64
	Action    Action
	Time      time.Time
	Amount    int
	BrockerID int
	ClientID  int
	OrderID   int32
	Ticker    string
	Partial   bool
	Instant   bool /* is from file ?*/
	Status    int32
}

/*
HasMorePriorityForSell - returns true if sell order a
has more priority than sell order b
*/
func HasMorePriorityForSell(a, b *Order) bool {
	return a.Price < b.Price || a.Time.Before(b.Time)
}

/*
HasMorePriorityForBuy - returns true if buy order a
has more priority than buy order b
*/
func HasMorePriorityForBuy(a, b *Order) bool {
	return a.Price > b.Price || a.Time.Before(b.Time)
}
