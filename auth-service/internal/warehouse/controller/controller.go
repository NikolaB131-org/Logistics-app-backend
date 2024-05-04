package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NikolaB131-org/logistics-backend/auth-service/internal/middlewares"
	"github.com/NikolaB131-org/logistics-backend/auth-service/otlp"
	proto "github.com/NikolaB131-org/logistics-backend/proto/warehouse"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	grpcClient proto.WarehouseServiceClient
}

func NewWarehouseRoutes(e *gin.Engine, grpcUrl string, middlewares middlewares.Middlewares) error {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("unable to connect warehouse grpc server: %w", err)
	}
	warehouseR := WarehouseRoutes{grpcClient: proto.NewWarehouseServiceClient(conn)}

	warehouse := e.Group("/warehouse", middlewares.OnlyAuth())
	{
		warehouse.GET("/products", warehouseR.handleGetProducts)
		warehouse.POST("/products", warehouseR.handleCreateProduct)
		warehouse.PATCH("/products/decrease", warehouseR.handleDecreaseProductQuantity, middlewares.OnlyWithRole("warehouseman"))
	}

	return nil
}

func (wc *WarehouseRoutes) handleGetProducts(c *gin.Context) {
	ctx, span := otlp.Tracer.Start(c, "HTTP GET /warehouse/products")
	defer span.End()

	res, err := wc.grpcClient.GetProducts(ctx, nil)
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	products := make([]GetProductsResponse, len(res.Products))
	for i, product := range res.Products {
		products[i] = GetProductsResponse{
			ID:       product.Id,
			Name:     product.Name,
			Quantity: int(product.Quantity),
			Price:    product.Price,
		}
	}

	c.JSON(http.StatusOK, products)
}

func (wc *WarehouseRoutes) handleCreateProduct(c *gin.Context) {
	ctx, span := otlp.Tracer.Start(c, "HTTP POST /warehouse/products")
	defer span.End()

	var req CreateProductRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	_, err := wc.grpcClient.CreateProduct(ctx, &proto.Product{Name: req.Name, Quantity: int64(req.Quantity), Price: req.Price})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (wc *WarehouseRoutes) handleDecreaseProductQuantity(c *gin.Context) {
	ctx, span := otlp.Tracer.Start(c, "HTTP PATCH /warehouse/products/decrease")
	defer span.End()

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

	_, err = wc.grpcClient.DecreaseProductQuantity(ctx, &proto.DecreaseProductQuantityRequest{Id: id, Quantity: int64(quantity)})
	if err != nil {
		slog.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
