package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	server, err := InitializeApp()

	if err != nil {
		log.Fatalf("Gagal Menginisialisasi aplikasi: %v", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"message": "Services is running",
		})
	})

	r.GET("/test-shutdown", func(c *gin.Context) {
		log.Println("[request started] memulai query/proses lambat...")

		var result int

		server.DB.Raw("SELECT pg_sleep(5)").Scan(&result)

		log.Println("[Request Finished] proses selesai")
		c.JSON(http.StatusOK, gin.H{"message": "Request berhasil diselesaikan dengan sukses"})
	})

	server.SetupRoutes(r)

	// log.Printf("Server berjalan di port %s", cfg.Port)
	// if err := r.Run(":" + cfg.Port); err != nil {
	// 	log.Fatalf("Gagal menjalankan server: %v", err)
	// }

	srv := &http.Server{
		Addr:    ":" + server.Config.Port,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Server berjalan di port: %s", server.Config.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Gagal menjalankan server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Signal shutdown diterima. Memulai graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Http server dihentikan paksa (timeout): %v", err)
	} else {
		log.Println("HTTP server berhasil dihentikan secara halus")
	}

	config.ClosePosgres(server.DB)

	log.Println("Aplikasi selesai dihentikan sepenuhnya")
}
