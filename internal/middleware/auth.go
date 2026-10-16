package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(cfg *config.AppConfig, rdb *redis.Client) gin.HandlerFunc {
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

		rawToken := parts[1]

		claims, err := helper.ValidateToken(parts[1], cfg.JWTSecret)
		if err != nil {
			helper.Error(c, http.StatusUnauthorized, "Token tidak valid atau telah kedaluwarsa")
			c.Abort()
			return
		}

		tokenHash := helper.HashTokenSHA256(rawToken)
		blacklistKey := fmt.Sprintf("blacklist:token:%s", tokenHash)

		isBlacklisted, err := rdb.Exists(c.Request.Context(), blacklistKey).Result()
		if err == nil && isBlacklisted > 0 {
			helper.Error(c, http.StatusUnauthorized, "sesi login telah berakhir, silahkan login kembali")
			c.Abort()
			return
		}

		c.Set("raw_token", rawToken)
		c.Set("token_claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
