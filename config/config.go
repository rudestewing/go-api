package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    AppPort     string
    JWTSecret   string
}

var GlobalConfig *Config

func InitConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: Error loading .env file:", err)
    }

    GlobalConfig = &Config{
        DatabaseURL: getEnv("DATABASE_URL", ""),
        AppPort:     getEnv("APP_PORT", "8000"),
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
    }

    log.Println("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func Get() *Config {
    if GlobalConfig == nil {
        log.Fatal("Config not initialized. Call InitConfig() first.")
    }
    return GlobalConfig
}
