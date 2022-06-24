# Inventory App
A practical web service for warehouses to manage inventory and orders.

## Features
- Create and manage products.
- Create, process, and store orders concurrently.
- Reject orders for out of stock products.
- Shutdown order intake channel.
- Load testing with an order simulator.

## Example Usages
#### 🚀 Start the service
```
go run server.go
```

#### 📦 Simulate multiple concurrent orders
```
go run simulations/loadtester.go
```

#### 🛒 View all products
```
curl localhost:8080/products
```

#### 💳 Place an order
```
curl -X POST -d '{"productId":"ARCEN", "quantity":5}' localhost:8080/orders/new
```

#### 🧾 View an order
```
curl localhost:8080/orders/<order_id>
```

#### 🚦 Stop incoming orders
```
curl localhost:8080/shutdown
```
