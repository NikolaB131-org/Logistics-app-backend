package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	authService "github.com/NikolaB131-org/logistics-backend/auth-service/internal/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	CheckToken(ctx context.Context, token string) bool
	ParseClaims(ctx context.Context, token string) (authService.TokenClaims, error)
}

type Middlewares struct {
	authService AuthService
}

func New(authService AuthService) Middlewares {
	return Middlewares{authService: authService}
}

func (m *Middlewares) OnlyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authorizationHeader, "Bearer ")
		token := splitToken[1]

		if len(token) < 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ok := m.authService.CheckToken(c, token); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ParseClaims(c, token)
		if err != nil {
			slog.Error(err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_roles", claims.RealmRoles)
		c.Next()
	}
}

func (m *Middlewares) OnlyWithRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("user_roles")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !slices.Contains(roles.([]string), role) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
