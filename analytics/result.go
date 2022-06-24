package analytics

import (
	"sync"

	"github.com/ndanny/inventory-app/models"
)

type result struct {
	data models.Analytics
	lock sync.Mutex
}

type Result interface {
	Get() models.Analytics
	Combine(analytics models.Analytics)
}

// Get returns the latest analytics data
func (r *result) Get() models.Analytics {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.data
}

// Combine updates the latest analytics data from incoming
// analytics events
func (r *result) Combine(analytics models.Analytics) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data = models.Combine(r.data, analytics)
}
