package main

import (
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
	authController "github.com/NikolaB131-org/logistics-backend/auth-service/internal/auth/controller"
	authService "github.com/NikolaB131-org/logistics-backend/auth-service/internal/auth/service"
	"github.com/NikolaB131-org/logistics-backend/auth-service/internal/middlewares"
	warehouseController "github.com/NikolaB131-org/logistics-backend/auth-service/internal/warehouse/controller"
	"github.com/NikolaB131-org/logistics-backend/auth-service/otlp"
	"github.com/gin-gonic/gin"
)

func main() {
	// OpenTelemetry
	otlp := otlp.New()
	otlp.Init()
	defer otlp.Shutdown()

	// Services
	authServ := authService.New(
		gocloak.NewClient(os.Getenv("KEYCLOAK_URL")),
		os.Getenv("KEYCLOAK_CLIENT_ID"),
		os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		os.Getenv("KEYCLOAK_REALM"),
	)

	// Middlewares
	middlewares := middlewares.New(authServ)

	// Routes
	r := gin.New()
	authController.NewAuthRoutes(r, authServ)
	err := warehouseController.NewWarehouseRoutes(r, os.Getenv("WAREHOUSE_GRPC_URL"), middlewares)
	if err != nil {
		panic(err)
	}

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
