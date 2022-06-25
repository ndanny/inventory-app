package db

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"

	"github.com/ndanny/inventory-app/models"
)

func ImportProducts(products *sync.Map) error {
	inp, err := readCsv()
	if err != nil {
		return err
	}

	for _, row := range inp {
		if len(row) != 4 {
			continue
		}
		var prod = models.Product{}
		prod.ID = row[0]
		prod.Name = row[1]
		if stock, err := strconv.Atoi(row[2]); err != nil {
			continue
		} else {
			prod.Stock = stock
		}
		if price, err := strconv.ParseFloat(row[3], 64); err != nil {
			continue
		} else {
			prod.Price = price
		}
		products.Store(prod.ID, prod)
	}

	return nil
}

func readCsv() ([][]string, error) {
	productsFile := "db/products.csv"
	f, err := os.Open(productsFile)
	defer f.Close()
	if err != nil {
		return [][]string{}, err
	}

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}
