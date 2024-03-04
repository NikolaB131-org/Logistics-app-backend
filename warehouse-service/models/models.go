package models

import (
	"context"
	"fmt"
	"os"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
)

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

func CreateProduct(product Product) {
	_, err := db.DbConn.Exec(context.Background(),
		`INSERT INTO products (name, quantity, price) VALUES($1, $2, $3)`,
		product.Name, product.Quantity, product.Price,
	)
	if err != nil {
		fmt.Println("Error while creating new product: ", err.Error())
		os.Exit(1)
	}
}

func GetProducts() []Product {
	rows, err := db.DbConn.Query(context.Background(), `SELECT * FROM products`)
	if err != nil {
		fmt.Println("Error while getting all products: ", err.Error())
		os.Exit(1)
	}

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
		products = append(products, product)
	}

	return products
}
