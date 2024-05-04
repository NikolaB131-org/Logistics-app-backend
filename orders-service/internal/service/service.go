package service

import (
	"context"

	"github.com/NikolaB131-org/logistics-backend/orders-service/internal/entity"
	proto "github.com/NikolaB131-org/logistics-backend/proto/warehouse"
)

type (
	OrdersRepository interface {
		CreateOrder(ctx context.Context, items []entity.Item) error
		GetOrders(ctx context.Context) ([]entity.Order, error)
	}

	Orders struct {
		ordersRepository OrdersRepository
		warehouseService proto.WarehouseServiceClient
	}
)

func NewOrdersService(ordersRepository OrdersRepository, warehouseService proto.WarehouseServiceClient) *Orders {
	return &Orders{ordersRepository: ordersRepository, warehouseService: warehouseService}
}

func (o *Orders) GetOrders(ctx context.Context) ([]entity.Order, error) {
	orders, err := o.ordersRepository.GetOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Orders) CreateOrder(ctx context.Context, items []entity.Item) error {
	for i := 0; i < len(items); i++ {
		o.warehouseService.DecreaseProductQuantity(ctx, &proto.DecreaseProductQuantityRequest{
			Id:       items[i].ID,
			Quantity: int64(items[i].Quantity),
		})
	}
	err := o.ordersRepository.CreateOrder(ctx, items)
	if err != nil {
		return err
	}
	return nil
}
