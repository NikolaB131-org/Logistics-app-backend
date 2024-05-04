package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
	grpccontroller "github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/controller/grpc"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/repository"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/service"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/rabbitmq"
	"google.golang.org/grpc"
)

func main() {
	err := db.ConnectDatabase(os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	rabbitmqClient, err := rabbitmq.New(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}

	productRepository := repository.NewProductRepository()

	warehouseService := service.NewWarehouseService(productRepository)

	gRPCServer := grpc.NewServer()

	grpccontroller.Register(gRPCServer, warehouseService, rabbitmqClient)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		panic(err)
	}
	slog.Info("gRPC server started", slog.String("port", os.Getenv("PORT")))

	if err := gRPCServer.Serve(listener); err != nil {
		panic(err)
	}
}
