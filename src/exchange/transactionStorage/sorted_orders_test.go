package transactionStorage

import (
	types "excange/transactionTypes"
	"testing"
	"time"
)

//TestSortedOrdersForSell - TODO
func TestSortedOrdersForSell(t *testing.T) {
	sortedOrders := NewSortedOrders(types.HasMorePriorityForSell)
	sortedOrders.Insert(
		&types.Order{
			Price: 100,
			Time:  time.Now(),
		})
}
