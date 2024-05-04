package main

import (
	"fmt"
	"os"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/controller"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/repository"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/service"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/rabbitmq"
	"github.com/gin-gonic/gin"
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

	// gRPCServer := grpc.NewServer()

	// grpccontroller.Register(gRPCServer, warehouseService, rabbitmqClient)

	// listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	// if err != nil {
	// 	panic(err)
	// }
	// slog.Info("gRPC server started", slog.String("port", os.Getenv("PORT")))

	// if err := gRPCServer.Serve(listener); err != nil {
	// 	panic(err)
	// }

	// Routes
	r := gin.New()

	err = controller.NewWarehouseRoutes(r, *warehouseService, *rabbitmqClient)
	if err != nil {
		panic(err)
	}

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
