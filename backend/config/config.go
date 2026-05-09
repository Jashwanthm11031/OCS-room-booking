package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
	Port          string
}

var AppConfig Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	AppConfig = Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "ocs_booking"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@ocs.iith.ac.in"),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
		Port:          getEnv("PORT", "8080"),
	}

	if AppConfig.DBPassword == "" {
		log.Fatal("DB_PASSWORD is required but not set")
	}
	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
