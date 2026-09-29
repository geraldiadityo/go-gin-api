package main

import (
	"log"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/modules/pegawai"
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

	api := r.Group("/api")
	pegawai.SetupRouter(api, pgDB)

	log.Printf("Server berjalan di port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
