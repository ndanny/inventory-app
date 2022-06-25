package main

import (
	"bytes"
	"encoding/json"
	"github.com/ndanny/inventory-app/models"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const concurrentOrders int = 50
const maxOrderQuantity int = 15
const newOrdersEndpoint string = "http://localhost:8080/orders/new"
const indexEndpoint string = "http://localhost:8080/"

var productIds []string = []string{"LBANE", "HGBLD", "ARCEN", "RACAP", "NEGCL", "HXREV", "NULLM", "MONSP"}

func main() {
	log.Println("Welcome to the concurrent load tester for the inventory app")
	log.Printf("Checking if %s is online...\n", indexEndpoint)
	// Exit if the service is not online.
	if err := serviceOnline(); err != nil {
		log.Fatalf("Error: service not online. Please start the server before this test.")
	}
	log.Printf("Good news! %s is online\n", indexEndpoint)
	log.Printf("Simulating %d concurrent orders on the service...\n", concurrentOrders)

	// Simulate multiple concurrent orders using goroutines
	// This will test for any possible race conditions, out-of-stock conditions, etc.
	var wg sync.WaitGroup
	wg.Add(concurrentOrders)
	for i := 0; i < concurrentOrders; i++ {
		go requestRngOrder(i, &wg)
	}
	wg.Wait()

	// Done
	log.Printf("Done!")
}

func requestRngOrder(num int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Mix up the rng seed
	rand.Seed(time.Now().UnixNano())

	// Create an order object
	quantity := rand.Intn(maxOrderQuantity) + 1
	pid := productIds[rand.Intn(len(productIds))]
	item := models.Item{
		ProductID: pid,
		Quantity:  quantity,
	}

	// Create http request and marshal the item to JSON
	ibytes, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", newOrdersEndpoint, bytes.NewBuffer(ibytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create the http client and send the request
	log.Printf("[sim-%d]: sending order %+v", num, item)
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[sim-%02d]: finished", num)
}

// Returns an error if the request to the service is faulty
func serviceOnline() error {
	req, err := http.NewRequest("GET", indexEndpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
