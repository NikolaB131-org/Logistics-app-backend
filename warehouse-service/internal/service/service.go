package service

import (
	"context"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/entity"
)

type (
	ProductRepository interface {
		CreateProduct(ctx context.Context, product entity.Product) error
		DecreaseProductQuantity(ctx context.Context, id string, quantity int) error
		GetProducts(ctx context.Context) ([]entity.Product, error)
	}

	Warehouse struct {
		productRepository ProductRepository
	}
)

func NewWarehouseService(productRepository ProductRepository) *Warehouse {
	return &Warehouse{productRepository: productRepository}
}

func (w *Warehouse) DecreaseProductQuantity(ctx context.Context, id string, quantity int) error {
	err := w.productRepository.DecreaseProductQuantity(ctx, id, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (w *Warehouse) GetProducts(ctx context.Context) ([]entity.Product, error) {
	products, err := w.productRepository.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (w *Warehouse) CreateProduct(ctx context.Context, product entity.Product) error {
	err := w.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return err
	}
	return nil
}
