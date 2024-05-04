package grpc

import (
	"context"
	"fmt"

	proto "github.com/NikolaB131-org/logistics-backend/proto/warehouse"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type WarehouseService interface {
	DecreaseProductQuantity(ctx context.Context, id string, quantity int) error
	GetProducts(ctx context.Context) ([]entity.Product, error)
	CreateProduct(ctx context.Context, product entity.Product) error
}

type RabbitMQ interface {
	SendNotification(message string)
}

type GRPCController struct {
	proto.UnimplementedWarehouseServiceServer
	warehouseService WarehouseService
	rabbitmq         RabbitMQ
}

func Register(gRPCServer *grpc.Server, warehouseService WarehouseService, rabbitmqClient RabbitMQ) {
	proto.RegisterWarehouseServiceServer(gRPCServer, &GRPCController{
		warehouseService: warehouseService,
		rabbitmq:         rabbitmqClient,
	})
}

func (c *GRPCController) DecreaseProductQuantity(ctx context.Context, req *proto.DecreaseProductQuantityRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	quantity := req.GetQuantity()

	err := c.warehouseService.DecreaseProductQuantity(ctx, id, int(quantity))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to decrease product quantity")
	}

	return nil, nil
}

func (c *GRPCController) GetProducts(ctx context.Context, req *emptypb.Empty) (*proto.GetProductsResponse, error) {
	products, err := c.warehouseService.GetProducts(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get products")
	}

	res := make([]*proto.Product, len(products))
	for i, product := range products {
		res[i] = &proto.Product{
			Id:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: int64(product.Quantity),
		}
	}

	return &proto.GetProductsResponse{Products: res}, nil
}

func (c *GRPCController) CreateProduct(ctx context.Context, req *proto.Product) (*emptypb.Empty, error) {
	name := req.GetName()
	price := req.GetPrice()
	quantity := req.GetQuantity()

	err := c.warehouseService.CreateProduct(ctx, entity.Product{Name: name, Price: price, Quantity: int(quantity)})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
	}
	c.rabbitmq.SendNotification(fmt.Sprintf("Product {name: %v, quantity: %v} was created", name, quantity))

	return nil, nil
}
