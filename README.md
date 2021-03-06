# Inventory App
A practical web service for warehouses to manage inventory and orders.

## Features
- Create and manage products.
- Create, process, and store orders concurrently.
- Reject orders for out of stock products.
- Shutdown order intake channel.
- Load testing with an order simulator.

## Example Usages
#### ๐ Start the service
```
go run -race server.go
```

#### ๐ฆ Simulate multiple concurrent orders
```
go run simulations/loadtester.go
```

#### ๐ View all products
```
curl localhost:8080/products
```

#### ๐ณ Place an order
```
curl -X POST -d '{"productId":"ARCEN", "quantity":5}' localhost:8080/orders/new
```

#### ๐งพ View an order
```
curl localhost:8080/orders/<order_id>
```

#### ๐ฆ Stop incoming orders
```
curl localhost:8080/shutdown
```
