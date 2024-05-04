package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/NikolaB131-org/logistics-backend/orders-service/db"
	grpccontroller "github.com/NikolaB131-org/logistics-backend/orders-service/internal/controller/grpc"
	"github.com/NikolaB131-org/logistics-backend/orders-service/internal/service"
	"github.com/NikolaB131-org/logistics-backend/orders-service/repository"
	warehouseProto "github.com/NikolaB131-org/logistics-backend/proto/warehouse"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := db.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(os.Getenv("WAREHOUSE_GRPC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("unable to connect grpc server: %w", err))
	}
	warehouseService := warehouseProto.NewWarehouseServiceClient(conn)

	ordersRepository := repository.NewOrdersRepository()

	ordersService := service.NewOrdersService(ordersRepository, warehouseService)

	gRPCServer := grpc.NewServer()

	grpccontroller.Register(gRPCServer, ordersService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}
	slog.Info("gRPC server started", slog.String("port", os.Getenv("PORT")))

	if err := gRPCServer.Serve(listener); err != nil {
		panic(err)
	}
}
