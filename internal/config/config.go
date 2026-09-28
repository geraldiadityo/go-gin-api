package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port  string
	PgDSN string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak di temukan menggunakan environtment sistem")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASSWORD", "secret")
	dbName := getEnv("DB_NAME", "db_name")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	dbTimezone := getEnv("DB_TIMEZONE", "Asia/Jakarta")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPass, dbName, dbPort, dbSSLMode, dbTimezone,
	)

	return &AppConfig{
		Port:  getEnv("PORT", "8080"),
		PgDSN: dsn,
	}
}
