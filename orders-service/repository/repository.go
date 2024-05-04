package repository

import (
	"context"
	"fmt"

	"github.com/NikolaB131-org/logistics-backend/orders-service/db"
	"github.com/NikolaB131-org/logistics-backend/orders-service/internal/entity"
	"github.com/jackc/pgx/v5"
)

type OrdersRepository struct{}

func NewOrdersRepository() *OrdersRepository {
	return &OrdersRepository{}
}

func (r *OrdersRepository) CreateOrder(ctx context.Context, items []entity.Item) error {
	_, err := db.DbConn.Exec(ctx, `INSERT INTO orders (items) VALUES($1)`, items)
	if err != nil {
		return fmt.Errorf("failed insert into: %w", err)
	}

	return nil
}

func (r *OrdersRepository) GetOrders(ctx context.Context) ([]entity.Order, error) {
	rows, err := db.DbConn.Query(ctx, "SELECT * FROM orders")
	if err != nil {
		return nil, fmt.Errorf("failed select: %w", err)
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("failed collecting rows: %w", err)
	}

	return orders, nil
}
