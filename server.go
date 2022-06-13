package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ndanny/inventory-app/handlers"
)

func main() {
	// Instantiate a handler
	handler, err := handlers.New()
	if err != nil {
		log.Fatalf("error creating new handler: %s", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/").Handler(http.HandlerFunc(handler.IndexHandler))
	router.Methods("GET").Path("/products").Handler(http.HandlerFunc(handler.ProductsHandler))
	router.Methods("GET").Path("/orders/{orderId}").Handler(http.HandlerFunc(handler.OrderGetHandler))
	router.Methods("GET").Path("/shutdown").Handler(http.HandlerFunc(handler.ShutdownHandler))
	router.Methods("POST").Path("/orders/new").Handler(http.HandlerFunc(handler.OrderCreateHandler))

	fmt.Println("Listening on localhost:8080...")
	http.ListenAndServe(":8080", router)
}
