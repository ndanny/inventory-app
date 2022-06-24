# Inventory App

## Examples
#### Start the service
```
go run server.go
```

#### Simulate multiple concurrent orders
```
go run simulations/loadtester.go
```

#### View all products
```
curl localhost:8080/products
```

#### Place an order
```
curl -X POST -d '{"productId":"ARCEN", "quantity":5}' localhost:8080/orders/new
```

#### View an order
```
curl localhost:8080/orders/<order_id>
```

#### Stop incoming orders
```
curl localhost:8080/shutdown
```
