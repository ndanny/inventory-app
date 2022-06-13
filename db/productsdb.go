package db

import (
	"fmt"

	"github.com/ndanny/inventory-app/models"
)

// ProductsDB simulates a database by storing kv pairs in
// an in-memory hashmap
type ProductsDB struct {
	products map[string]models.Product
}

func NewProductsDB() (*ProductsDB, error) {
	p := &ProductsDB{
		products: make(map[string]models.Product),
	}
	if err := ImportProducts(p.products); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *ProductsDB) Exists(id string) error {
	if _, ok := p.products[id]; !ok {
		return fmt.Errorf("product id %s not found in the db", id)
	}

	return nil
}

func (p *ProductsDB) Find(id string) (models.Product, error) {
	prod, ok := p.products[id]
	if !ok {
		return models.Product{}, fmt.Errorf("product id %s not found in the db", id)
	}

	return prod, nil
}

func (p *ProductsDB) Insert(product models.Product) {
	p.products[product.ID] = product
}

func (p *ProductsDB) GetAll() []models.Product {
	all := make([]models.Product, 0, len(p.products))
	for _, prod := range p.products {
		all = append(all, prod)
	}

	return all
}
