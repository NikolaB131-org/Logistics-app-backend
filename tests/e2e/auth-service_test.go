package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AuthServiceSuite struct {
	suite.Suite

	httpClient *http.Client
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceSuite))
}

func (s *AuthServiceSuite) SetupSuite() {
	s.httpClient = &http.Client{Timeout: 10 * time.Second}
}

func (s *AuthServiceSuite) TestWithoutRoles() {
	body, _ := json.Marshal(map[string]string{
		"email":    "testemail@go.com",
		"password": "qwerty",
	})

	// Register
	res, err := http.Post("http://localhost:4000/auth/register", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.NotEmpty(data["userId"])

	// Login
	res, err = http.Post("http://localhost:4000/auth/login", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.NotEmpty(data["token"])
}

func (s *AuthServiceSuite) TestWithWarehouseRole() {
	body, _ := json.Marshal(map[string]string{
		"email":    "testemail2@go.com",
		"password": "qwerty",
		"role":     "warehouseman",
	})

	// Register
	res, err := http.Post("http://localhost:4000/auth/register", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.NotEmpty(data["userId"])

	// Login
	res, err = http.Post("http://localhost:4000/auth/login", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.NotEmpty(data["token"])
}

func (s *AuthServiceSuite) TestWithIncorrectRole() {
	body, _ := json.Marshal(map[string]string{
		"email":    "testemail2@go.com",
		"password": "qwerty",
		"role":     "warehouseman123",
	})

	// Register
	res, err := http.Post("http://localhost:4000/auth/register", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.Equal(data["error"], "invalid role name")
}

func (s *AuthServiceSuite) TestToken() {
	body, _ := json.Marshal(map[string]string{
		"email":    "testemail3@go.com",
		"password": "qwerty",
	})

	// Register
	res, err := http.Post("http://localhost:4000/auth/register", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.NotEmpty(data["userId"])

	// Login
	res, err = http.Post("http://localhost:4000/auth/login", "application/json", bytes.NewBuffer(body))
	s.Require().NoError(err)

	err = json.NewDecoder(res.Body).Decode(&data)
	s.Require().NoError(err)
	s.NotEmpty(data["token"])
	token := data["token"]

	// Valid token
	req, err := http.NewRequest(http.MethodGet, "http://localhost:4000/warehouse/products", nil)
	s.Require().NoError(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err = s.httpClient.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)

	// Invalid token
	req, err = http.NewRequest(http.MethodGet, "http://localhost:4000/warehouse/products", nil)
	s.Require().NoError(err)
	req.Header.Add("Authorization", "Bearer wrong_token")

	res, err = s.httpClient.Do(req)
	s.Require().NoError(err)
	s.Equal(http.StatusUnauthorized, res.StatusCode)
}
