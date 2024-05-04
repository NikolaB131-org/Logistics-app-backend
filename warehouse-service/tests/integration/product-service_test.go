package integration

import (
	"context"
	"slices"
	"testing"

	"github.com/NikolaB131-org/logistics-backend/warehouse-service/db"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/entity"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/repository"
	"github.com/NikolaB131-org/logistics-backend/warehouse-service/internal/service"
	"github.com/stretchr/testify/suite"
)

type ProductServiceSuite struct {
	suite.Suite

	service *service.Warehouse
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}

func (s *ProductServiceSuite) SetupSuite() {
	err := db.ConnectDatabase("postgresql://postgres:postgres@localhost:7777/warehouse?sslmode=disable")
	s.Require().NoError(err)

	productRepository := repository.NewProductRepository()
	s.service = service.NewWarehouseService(productRepository)
}

func (s *ProductServiceSuite) TestCreateProduct() {
	ctx := context.Background()
	testProduct := entity.Product{Name: "test_name", Quantity: 56, Price: 2.49}

	err := s.service.CreateProduct(ctx, testProduct)
	s.Require().NoError(err)

	var resProduct entity.Product
	err = db.DbConn.QueryRow(ctx, "SELECT * FROM products WHERE name='test_name' AND quantity = 56").Scan(&resProduct.ID, &resProduct.Name, &resProduct.Quantity, &resProduct.Price)
	s.Require().NoError(err)

	s.Equal(testProduct.Name, resProduct.Name)
	s.Equal(testProduct.Quantity, resProduct.Quantity)
	s.Equal(testProduct.Price, resProduct.Price)
}

func (s *ProductServiceSuite) TestGetProducts() {
	ctx := context.Background()
	testProduct := entity.Product{Name: "test_name2", Quantity: 58, Price: 2.66}
	testProduct2 := entity.Product{Name: "test_name3", Quantity: 2, Price: 6.68}

	err := s.service.CreateProduct(ctx, testProduct)
	s.Require().NoError(err)
	err = s.service.CreateProduct(ctx, testProduct2)
	s.Require().NoError(err)

	products, err := s.service.GetProducts(ctx)
	s.Require().NoError(err)

	idx1 := slices.IndexFunc(products, func(p entity.Product) bool { return p.Name == "test_name2" })
	idx2 := slices.IndexFunc(products, func(p entity.Product) bool { return p.Name == "test_name3" })

	s.NotEqual(-1, idx1)
	s.NotEqual(-1, idx2)

	s.Equal(testProduct.Name, products[idx1].Name)
	s.Equal(testProduct.Quantity, products[idx1].Quantity)
	s.Equal(testProduct.Price, products[idx1].Price)

	s.Equal(testProduct2.Name, products[idx2].Name)
	s.Equal(testProduct2.Quantity, products[idx2].Quantity)
	s.Equal(testProduct2.Price, products[idx2].Price)
}

func (s *ProductServiceSuite) TestDecreaseProductQuantity() {
	ctx := context.Background()
	testProduct := entity.Product{Name: "test_name4", Quantity: 111, Price: 1.66}

	err := s.service.CreateProduct(ctx, testProduct)
	s.Require().NoError(err)

	var id string
	err = db.DbConn.QueryRow(ctx, "SELECT id FROM products WHERE name='test_name4'").Scan(&id)
	s.Require().NoError(err)

	err = s.service.DecreaseProductQuantity(ctx, id, 2)
	s.Require().NoError(err)

	var quantity int
	err = db.DbConn.QueryRow(ctx, "SELECT quantity FROM products WHERE name='test_name4'").Scan(&quantity)
	s.Require().NoError(err)

	s.Equal(109, quantity)

	err = s.service.DecreaseProductQuantity(ctx, id, 110)
	s.Equal("error while updating product: future quantity less than 0", err.Error())
}
