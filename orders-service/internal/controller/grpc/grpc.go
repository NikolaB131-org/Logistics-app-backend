package grpc

import (
	"context"
	"log/slog"

	"github.com/NikolaB131-org/logistics-backend/orders-service/internal/entity"
	proto "github.com/NikolaB131-org/logistics-backend/proto/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrdersService interface {
	GetOrders(ctx context.Context) ([]entity.Order, error)
	CreateOrder(ctx context.Context, items []entity.Item) error
}

type GRPCController struct {
	proto.UnimplementedOrdersServiceServer
	ordersService OrdersService
}

func Register(gRPCServer *grpc.Server, ordersService OrdersService) {
	proto.RegisterOrdersServiceServer(gRPCServer, &GRPCController{
		ordersService: ordersService,
	})
}

func (c *GRPCController) GetOrders(ctx context.Context, req *emptypb.Empty) (*proto.GetOrdersResponse, error) {
	orders, err := c.ordersService.GetOrders(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to get orders")
	}

	convertedOrders := make([]*proto.Order, len(orders))
	for i, order := range orders {
		items := make([]*proto.Item, len(order.Items))
		for j, product := range order.Items {
			items[j] = &proto.Item{Id: product.ID, Quantity: int64(product.Quantity)}
		}
		convertedOrders[i] = &proto.Order{Items: items}
	}

	return &proto.GetOrdersResponse{Orders: convertedOrders}, nil
}

func (c *GRPCController) CreateOrder(ctx context.Context, req *proto.Order) (*emptypb.Empty, error) {
	items := req.GetItems()

	convertedItems := make([]entity.Item, len(items))
	for i, item := range items {
		convertedItems[i] = entity.Item{ID: item.Id, Quantity: int(item.Quantity)}
	}

	err := c.ordersService.CreateOrder(ctx, convertedItems)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Error(codes.Internal, "failed to create order")
	}

	return nil, nil
}
