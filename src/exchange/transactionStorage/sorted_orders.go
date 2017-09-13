package transactionStorage

import (
	_ "container/heap"
	types "exchange/transactionTypes"
	"fmt"
	"sort"
)

var (
	errEmptyStorage = fmt.Errorf("Empty storage")
	errNoSuchItem   = fmt.Errorf("No such item")
)

//SortedOrdersInterface -
type SortedOrdersInterface interface {
	Insert(*types.Order)
	Peek() (*types.Order, error)
	Pop() (*types.Order, error)
	Delete(BrokerID int, OrderID int32) error
}

//SortedOrders - сортировка по возрастанию
type SortedOrders struct {
	SortedOrdersInterface
	hasMorePriority func(a, b *types.Order) bool
	orders          []*types.Order
}

/*
NewSortedOrders -
hasMorePriority returns true if a has more priority than b
*/
func NewSortedOrders(hasMorePriority func(a, b *types.Order) bool) *SortedOrders {
	return &SortedOrders{
		hasMorePriority: hasMorePriority,
		orders:          make([]*types.Order, 0, 100),
	}
}

//Insert - inserts a new order
func (s *SortedOrders) Insert(o *types.Order) {
	s.orders = append(s.orders, o)
	sort.Slice(s.orders, func(i, j int) bool {
		return s.hasMorePriority(s.orders[i], s.orders[j])
	})
}

//Peek - returns order with highest priority
func (s *SortedOrders) Peek() (*types.Order, error) {
	if len(s.orders) == 0 {
		return nil, errEmptyStorage
	}
	return s.orders[0], nil
}

//Pop - extracts and returns order with highest priority
func (s *SortedOrders) Pop() (*types.Order, error) {
	if len(s.orders) == 0 {
		return nil, errEmptyStorage
	}
	o := s.orders[0]
	s.orders = s.orders[1:]
	return o, nil
}

func (s *SortedOrders) Delete(BrockerID int, OrderID int32) error {
	if len(s.orders) == 0 {
		return errEmptyStorage
	}
	length := len(s.orders)
	for i := 0; i < length; i++ {
		if s.orders[i].BrockerID != BrockerID || s.orders[i].OrderID != OrderID {
			continue
		}
		if i == length-1 {
			s.orders = s.orders[:i]
		} else {
			s.orders = append(s.orders[:i], s.orders[i+1:]...)
		}
		return nil
	}

	return errNoSuchItem
}

/* heap interface */
type ordersHeap struct {
	orders          []*types.Order
	hasMorePriority func(a, b *types.Order) bool
}

func (h *ordersHeap) Len() int {
	return len(h.orders)
}

func (h *ordersHeap) Less(i, j int) bool {
	return h.hasMorePriority(h.orders[i], h.orders[j])
}

func (h *ordersHeap) Swap(i, j int) {
	h.orders[i], h.orders[j] = h.orders[j], h.orders[i]
}

func (h *ordersHeap) Push(x interface{}) {

}
