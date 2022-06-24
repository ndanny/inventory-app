package analytics

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ndanny/inventory-app/models"
)

const WorkerCount = 3

type analyticsService struct {
	result     Result
	processed  <-chan models.Order
	done       <-chan struct{}
	pAnalytics chan models.Analytics
}

type AnalyticsService interface {
	GetAnalytics() models.Analytics
}

// New returns a new instance of an AnalyticsService
func New(processed <-chan models.Order, done <-chan struct{}) AnalyticsService {
	a := analyticsService{
		result:     &result{},
		processed:  processed,
		done:       done,
		pAnalytics: make(chan models.Analytics, WorkerCount),
	}
	// Create worker pool to process multiple orders at once
	for i := 0; i < WorkerCount; i++ {
		go a.run()
	}
	go a.reconcile()
	return &a
}

// GetAnalytics returns the latest analytics data for the inventory app
func (a *analyticsService) GetAnalytics() models.Analytics {
	return a.result.Get()
}

// run listens to incoming orders to update the overall analytics
func (a *analyticsService) run() {
	fmt.Println("Gathering service analytics...")
	for {
		select {
		case order := <-a.processed:
			// Updates the service analytics from incoming order events
			p := a.processOrder(order)
			a.pAnalytics <- p
		case <-a.done:
			fmt.Println("Analytics service has stopped")
			return
		}
	}
}

// reconcile listens to the a.pAnalytics chan for analytics events and
// combines the data with the latest analytics data collected
func (a *analyticsService) reconcile() {
	fmt.Println("Reconcile started")
	for {
		select {
		case p := <-a.pAnalytics:
			a.result.Combine(p)
		case <-a.done:
			fmt.Println("Reconcile stopped")
			return
		}
	}
}

// processOrder takes an order and returns an analytics event based on
// the order completion status
func (a *analyticsService) processOrder(order models.Order) models.Analytics {
	// Simulate a costly operation
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

	data := models.Analytics{}
	if order.Status == models.OrderStatusCompleted {
		data.CompletedOrders = 1
		data.TotalRevenue = order.Total
	} else {
		data.RejectedOrders = 1
	}

	return data
}
