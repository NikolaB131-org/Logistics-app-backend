package models

import (
	"context"
	"errors"
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

type DecreaseProductQuantityBody struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
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

func DecreaseProductQuantity(b DecreaseProductQuantityBody) error {
	rows, err := db.DbConn.Query(context.Background(), `SELECT * FROM products WHERE id = $1`, b.Id)
	if err != nil {
		return fmt.Errorf("error while getting product: %s", err.Error())
	}
	var product Product
	for rows.Next() {
		rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
	}

	if product.Quantity-b.Quantity < 0 {
		return errors.New("error while updating product: future quantity less than 0")
	}

	_, err = db.DbConn.Exec(context.Background(), "UPDATE products SET quantity = $1 WHERE id = $2", product.Quantity-b.Quantity, b.Id)
	if err != nil {
		return fmt.Errorf("error while updating product: %s", err.Error())
	}

	return nil
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
