package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/geraldiadityo/go-backend/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func LoginRateLimiter(rdb *redis.Client, maxRequest int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:login:ip:%s", ip)
		ctx := c.Request.Context()
		count, err := rdb.Incr(ctx, key).Result()

		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > maxRequest {
			ttl, _ := rdb.TTL(ctx, key).Result()
			c.Header("Retry-After", fmt.Sprintf("%.0f", ttl.Seconds()))
			helper.Error(c, http.StatusTooManyRequests, fmt.Sprintf("Terlalu banyak percobaan login dari IP ini. coba lagi dalam %.0f detik", ttl.Seconds()))
			c.Abort()
			return
		}

		c.Next()
	}
}
