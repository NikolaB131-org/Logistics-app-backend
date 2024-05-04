package controller

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/entity"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/service"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/rabbitmq"
	"github.com/gin-gonic/gin"
)

type CreateProductRequest struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

type GetProductsResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"price"`
}

type WarehouseRoutes struct {
	// grpcClient proto.WarehouseServiceClient
	warehouseService service.Warehouse
	rabbitmq         rabbitmq.RabbitMQ
}

func NewWarehouseRoutes(e *gin.Engine, warehouseService service.Warehouse, rabbitmq rabbitmq.RabbitMQ) error {
	// conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	return fmt.Errorf("unable to connect warehouse grpc server: %w", err)
	// }
	warehouseR := WarehouseRoutes{warehouseService: warehouseService, rabbitmq: rabbitmq}

	warehouse := e.Group("/warehouse")
	{
		warehouse.GET("/products", warehouseR.handleGetProducts)
		warehouse.POST("/products", warehouseR.handleCreateProduct)
		warehouse.PATCH("/products/decrease", warehouseR.handleDecreaseProductQuantity)
	}

	return nil
}

func (wc *WarehouseRoutes) handleGetProducts(c *gin.Context) {
	// ctx, span := otlp.Tracer.Start(c, "HTTP GET /warehouse/products")
	// defer span.End()

	products, err := wc.warehouseService.GetProducts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}

	// res, err := wc.grpcClient.GetProducts(ctx, nil)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	c.Status(http.StatusInternalServerError)
	// 	return
	// }
	// products = make([]GetProductsResponse, len(res.Products))
	// for i, product := range res.Products {
	// 	products[i] = GetProductsResponse{
	// 		ID:       product.Id,
	// 		Name:     product.Name,
	// 		Quantity: int(product.Quantity),
	// 		Price:    product.Price,
	// 	}
	// }

	c.JSON(http.StatusOK, products)
}

func (wc *WarehouseRoutes) handleCreateProduct(c *gin.Context) {
	// ctx, span := otlp.Tracer.Start(c, "HTTP POST /warehouse/products")
	// defer span.End()

	var req CreateProductRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	err := wc.warehouseService.CreateProduct(c, entity.Product{Name: req.Name, Quantity: req.Quantity, Price: req.Price})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (wc *WarehouseRoutes) handleDecreaseProductQuantity(c *gin.Context) {
	// ctx, span := otlp.Tracer.Start(c, "HTTP PATCH /warehouse/products/decrease")
	// defer span.End()

	id := c.Query("id")
	quantityQ := c.Query("quantity")
	quantity, err := strconv.Atoi(quantityQ)
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	if quantity < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity is required and must be > 0"})
		return
	}

	err = wc.warehouseService.DecreaseProductQuantity(c, id, quantity)
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
