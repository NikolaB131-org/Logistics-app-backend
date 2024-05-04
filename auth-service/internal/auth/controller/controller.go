package controller

import (
	"context"
	"errors"
	"log"
	"net/http"

	authService "github.com/NikolaB131-org/logistics-backend/auth-service/internal/auth/service"
	"github.com/NikolaB131-org/logistics-backend/auth-service/otlp"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string, role string) (string, error)
}

type AuthRoutes struct {
	authService AuthService
}

func NewAuthRoutes(e *gin.Engine, authService AuthService) {
	authR := AuthRoutes{authService: authService}

	auth := e.Group("/auth")
	{
		auth.POST("/login", authR.handleLogin)
		auth.POST("/register", authR.handleRegister)
	}
}

func (r *AuthRoutes) handleLogin(c *gin.Context) {
	ctx, span := otlp.Tracer.Start(c, "HTTP POST /login")
	defer span.End()

	var req LoginRequest

	if err := c.BindJSON(&req); err != nil {
		log.Print(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	token, err := r.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *AuthRoutes) handleRegister(c *gin.Context) {
	ctx, span := otlp.Tracer.Start(c, "HTTP POST /register")
	defer span.End()

	var req RegisterRequest

	if err := c.BindJSON(&req); err != nil {
		log.Print(err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	userId, err := r.authService.Register(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidRoleName) || errors.Is(err, authService.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": userId})
}
