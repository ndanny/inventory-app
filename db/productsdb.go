package db

import (
	"fmt"
	"sort"
	"sync"

	"github.com/ndanny/inventory-app/models"
)

// ProductsDB simulates a database by storing kv pairs in
// an in-memory hashmap
type ProductsDB struct {
	products sync.Map
}

func NewProductsDB() (*ProductsDB, error) {
	p := &ProductsDB{}
	if err := ImportProducts(&p.products); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *ProductsDB) Exists(id string) error {
	if _, ok := p.products.Load(id); !ok {
		return fmt.Errorf("product id %s not found in the db", id)
	}

	return nil
}

func (p *ProductsDB) Find(id string) (models.Product, error) {
	prod, ok := p.products.Load(id)
	if !ok {
		return models.Product{}, fmt.Errorf("product id %s not found in the db", id)
	}

	return p.toProduct(prod), nil
}

func (p *ProductsDB) Insert(product models.Product) {
	p.products.Store(product.ID, product)
}

func (p *ProductsDB) GetAll() []models.Product {
	var everything []models.Product
	p.products.Range(func(_, pp interface{}) bool {
		everything = append(everything, p.toProduct(pp))
		return true
	})
	sort.Slice(everything, func(i, j int) bool {
		return everything[i].ID < everything[j].ID
	})

	return everything
}

func (p *ProductsDB) toProduct(pp interface{}) models.Product {
	prod, ok := pp.(models.Product)
	if !ok {
		panic(fmt.Errorf("error casting %v to product", pp))
	}

	return prod
}
