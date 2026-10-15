package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *AppConfig) *redis.Client {
	opt, err := redis.ParseURL(cfg.RedisDSN)
	if err != nil {
		log.Fatalf("Format REDIS_DSN tidak valid: %v", err)
	}

	opt.PoolSize = 20
	opt.MinIdleConns = 5

	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Peringatan: Gagal terhubung ke redis (%v). Aplikasi tetap berjalan (fallback ke DB)", err)
	} else {
		log.Println("Redis terkoneksi dengan connection pool aktif")
	}

	return rdb
}

func CloseRedis(rdb *redis.Client) {
	if rdb == nil {
		return
	}

	if err := rdb.Close(); err != nil {
		log.Printf("Gagal menutup koneksi redis: %v", err)
	} else {
		log.Println("Koneksi redis berhasil di tutup bersih")
	}
}
