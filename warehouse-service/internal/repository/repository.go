package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/entity"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product entity.Product) error {
	_, err := db.DbConn.Exec(ctx,
		`INSERT INTO products (name, quantity, price) VALUES($1, $2, $3)`,
		product.Name, product.Quantity, product.Price,
	)
	if err != nil {
		return fmt.Errorf("failed insert into: %w", err)
	}
	return nil
}

func (r *ProductRepository) DecreaseProductQuantity(ctx context.Context, id string, quantity int) error {
	rows, err := db.DbConn.Query(ctx, `SELECT * FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed select: %w", err)
	}
	var product entity.Product
	for rows.Next() {
		rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
	}

	if product.Quantity-quantity < 0 {
		return errors.New("error while updating product: future quantity less than 0")
	}

	_, err = db.DbConn.Exec(ctx, "UPDATE products SET quantity = $1 WHERE id = $2", product.Quantity-quantity, id)
	if err != nil {
		return fmt.Errorf("failed update: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetProducts(ctx context.Context) ([]entity.Product, error) {
	rows, err := db.DbConn.Query(ctx, `SELECT * FROM products`)
	if err != nil {
		return []entity.Product{}, fmt.Errorf("failed select: %w", err)
	}

	products := make([]entity.Product, 0)
	for rows.Next() {
		var product entity.Product
		rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
		products = append(products, product)
	}

	return products, nil
}
