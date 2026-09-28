package main

import (
	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	pgDB := config.InitPostgres(cfg.PgDSN)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"message": "Services is running",
		})
	})
}
