package middleware

import (
	"net/http"
	"strings"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.Error(c, http.StatusUnauthorized, "Authorization header tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helper.Error(c, http.StatusUnauthorized, "Format token harus 'Bearer <token>'")
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(parts[1], cfg.JWTSecret)
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, "Token tidak valid atau telah kedaluwarsa")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

