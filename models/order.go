package models

import (
	"github.com/google/uuid"
	"time"
)

// OrderStatus defines enums for order statuses
type OrderStatus string

const (
	OrderStatusNew       OrderStatus = "New"
	OrderStatusRejected  OrderStatus = "Rejected"
	OrderStatusCompleted OrderStatus = "Completed"
)

type Order struct {
	ID        string      `json:"id"`
	Item      Item        `json:"item"`
	Total     float64     `json:"total"`
	Status    OrderStatus `json:"status"`
	Error     string      `json:"error"`
	CreatedAt string      `json:"createdAt"`
}

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func NewOrder(item Item) Order {
	return Order{
		ID:        uuid.NewString(),
		Item:      item,
		Status:    OrderStatusNew,
		CreatedAt: time.Now().String(),
	}
}

func (o *Order) Done() {
	o.Status = OrderStatusCompleted
}
