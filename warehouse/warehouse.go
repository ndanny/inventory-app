package warehouse

import (
	"fmt"
	"math"

	"github.com/ndanny/inventory-app/db"
	"github.com/ndanny/inventory-app/models"
)

// warehouse manages an inventory of products and orders
type warehouse struct {
	products *db.ProductsDB
	orders   *db.OrdersDB
	incoming chan models.Order
	done     chan struct{}
}

type Warehouse interface {
	GetProducts() []models.Product
	GetOrder(id string) (models.Order, error)
	CreateOrder(item models.Item) (*models.Order, error)
	Close()
}

func New() (Warehouse, error) {
	p, err := db.NewProductsDB()
	if err != nil {
		return nil, err
	}
	o, err := db.NewOrdersDB()
	if err != nil {
		return nil, err
	}
	wh := warehouse{
		products: p,
		orders:   o,
		incoming: make(chan models.Order),
		done:     make(chan struct{}),
	}

	// When creating a new warehouse, spin up a goroutine to continually
	// listen for incoming orders on the incoming channel.
	go wh.processOrders()

	return &wh, nil
}

func (w *warehouse) GetProducts() []models.Product {
	return w.products.GetAll()
}

func (w *warehouse) GetOrder(id string) (models.Order, error) {
	return w.orders.FindOrder(id)
}

func (w *warehouse) CreateOrder(item models.Item) (*models.Order, error) {
	err := w.validateItem(item)
	if err != nil {
		return nil, err
	}
	o := models.NewOrder(item)

	// Upon processing a new order, there are two cases:
	// 1. Accept the new order and place the order in a channel for incoming orders
	// 2. The done channel is closed, so do not accept the new order
	select {
	case w.incoming <- o:
		w.orders.Insert(o)
		fmt.Printf("Order created for product id %s\n", item.ProductID)
		return &o, nil
	case <-w.done:
		return nil, fmt.Errorf("not taking orders at the moment")
	}
}

func (w *warehouse) Close() {
	// Closes the done channel, indicating that no more orders can be received
	close(w.done)
}

func (w *warehouse) validateItem(item models.Item) error {
	if item.Quantity < 1 {
		return fmt.Errorf("quantity must be at least one: got %d", item.Quantity)
	}
	if err := w.products.Exists(item.ProductID); err != nil {
		return fmt.Errorf("product id %s does not exist", item.ProductID)
	}

	return nil
}

// processOrders listens for new orders on the incoming channel
func (w *warehouse) processOrders() {
	fmt.Println("Order processor started! Listening for orders...")
	for {
		select {
		case order := <-w.incoming:
			fmt.Printf("Processing order %s...\n", order.ID)
			w.processOrder(&order)
			w.orders.Insert(order)
		case <-w.done:
			fmt.Println("Order processor shutting down...")
			return
		}
	}
}

// processOrder completes or rejects an order based on order conditions
func (w *warehouse) processOrder(order *models.Order) {
	item := order.Item
	prod, err := w.products.Find(item.ProductID)
	// Reject order if product can not be found
	if err != nil {
		order.Status = models.OrderStatusRejected
		order.Error = err.Error()
		return
	}
	// Reject order if stock on hand is less than requested quantity
	if prod.Stock < item.Quantity {
		order.Status = models.OrderStatusRejected
		order.Error = fmt.Sprintf("not enough stock for %s with stock %d: got %d",
			item.ProductID, prod.Stock, item.Quantity)
		return
	}
	// Update the product in the products database
	prod.Stock -= item.Quantity
	w.products.Insert(prod)
	// Update the order with the total cost and mark it as done
	order.Total = math.Round(float64(order.Item.Quantity)*prod.Price*100) / 100
	order.Done()
}
