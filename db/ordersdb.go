package db

import (
	"fmt"
	"sync"

	"github.com/ndanny/inventory-app/models"
)

type OrdersDB struct {
	// https://pkg.go.dev/sync#Map
	orderRequests sync.Map
}

func NewOrdersDB() (*OrdersDB, error) {
	return &OrdersDB{}, nil
}

func (o *OrdersDB) FindOrder(id string) (models.Order, error) {
	order, ok := o.orderRequests.Load(id)
	if !ok {
		return models.Order{}, fmt.Errorf("no order found for id %s", id)
	}

	// Casts the order interface into an internal order model
	conv, ok := order.(models.Order)
	if !ok {
		return models.Order{}, fmt.Errorf("cannot cast %v to order", order)
	}

	return conv, nil
}

// Insert creates or updates an order in the orderRequests map
func (o *OrdersDB) Insert(order models.Order) {
	o.orderRequests.Store(order.ID, order)
}
