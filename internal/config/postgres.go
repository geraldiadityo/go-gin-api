package config

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Gagal terkoneksi ke postgre SQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal menginitialisasi connection pool: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("PostgreSQL terkoneksi dengan connection pool aktif")

	return db
}

func ClosePosgres(db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Gagal mendapatkan sql.DB saat akan menutup koneksi: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Gagal menutup koneksi PostgreSQL: %v", err)
	} else {
		log.Println("koneksi PosgreSQL berhasil di tutup bersih")
	}
}

func ProvideDB(cfg *AppConfig) *gorm.DB {
	return InitPostgres(cfg.PgDSN)
}
