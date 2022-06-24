package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ndanny/inventory-app/models"
	"github.com/ndanny/inventory-app/warehouse"
)

type handler struct {
	wh   warehouse.Warehouse
	once sync.Once
}

type Handler interface {
	IndexHandler(w http.ResponseWriter, r *http.Request)
	ProductsHandler(w http.ResponseWriter, r *http.Request)
	OrderGetHandler(w http.ResponseWriter, r *http.Request)
	OrderCreateHandler(w http.ResponseWriter, r *http.Request)
	ShutdownHandler(w http.ResponseWriter, r *http.Request)
}

// New returns an instance of the handler object
func New() (Handler, error) {
	wh, err := warehouse.New()
	if err != nil {
		return nil, err
	}
	h := handler{
		wh: wh,
	}
	return &h, nil
}

// IndexHandler returns a simple string response upon visiting the home page
func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Welcome to the inventory app!"
	writeResponse(w, http.StatusOK, msg, nil)
}

// ProductsHandler shows all products in the database
func (h *handler) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	prods := h.wh.GetProducts()
	writeResponse(w, http.StatusOK, prods, nil)
}

// OrderGetHandler fetches an order from the OrdersDB and returns it
func (h *handler) OrderGetHandler(w http.ResponseWriter, r *http.Request) {
	// mux.Vars() gets a map of route variables from the URL path ("/orders/{orderId}")
	vars := mux.Vars(r)
	id := vars["orderId"]

	// Look into the warehouse OrdersDB for the id
	o, err := h.wh.GetOrder(id)
	if err != nil {
		writeResponse(w, http.StatusNotFound, nil, err)
		return
	}

	writeResponse(w, http.StatusOK, o, nil)
}

// OrderCreateHandler creates a new order from the requested order
func (h *handler) OrderCreateHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&item); err != nil {
		writeResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid order request body"))
		return
	}

	o, err := h.wh.CreateOrder(item)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	writeResponse(w, http.StatusOK, o, nil)
}

// ShutdownHandler closes warehouse from taking new orders
func (h *handler) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Closing the warehouse from taking new orders...")
	// sync.Once ensures that the given method is only invoked once
	h.once.Do(func() {
		h.wh.Close()
	})
	writeResponse(w, http.StatusOK, "The warehouse is now closed!", nil)
}
