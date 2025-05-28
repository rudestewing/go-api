package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	AppPort     string
	JWTSecret   string
	// Security configurations
	JWTExpiry   time.Duration
	Environment string
	// Database configurations
	DBMaxIdleConns int
	DBMaxOpenConns int
	DBMaxLifetime  time.Duration
	// Server timeout configurations
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	// Logging configurations
	LogDir         string
	LogMaxSize     int64
	LogMaxAge      int
	EnableDailyLog bool
}

var GlobalConfig *Config

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	validateRequiredEnvs()

	GlobalConfig = &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AppPort:     getEnv("APP_PORT", "8000"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		// Security configurations
		JWTExpiry:   getEnvAsDuration("JWT_EXPIRY", time.Hour*24),
		Environment: getEnv("ENVIRONMENT", "development"),
		// Database configurations
		DBMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		DBMaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", time.Hour),
		// Server timeout configurations
		ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", time.Second*30),
		WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", time.Second*30),
		IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", time.Second*120),
		// Logging configurations
		LogDir:         getEnv("LOG_DIR", "storage/logs"),
		LogMaxSize:     getEnvAsInt64("LOG_MAX_SIZE", 10*1024*1024), // 10MB
		LogMaxAge:      getEnvAsInt("LOG_MAX_AGE", 30),              // 30 days
		EnableDailyLog: getEnvAsBool("ENABLE_DAILY_LOG", true),
	}

	log.Println("Configuration loaded successfully")
}

// validateRequiredEnvs checks if all required environment variables are set
func validateRequiredEnvs() {
	requiredEnvs := []string{
		"DATABASE_URL",
		"JWT_SECRET",
	}

	var missingEnvs []string
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			missingEnvs = append(missingEnvs, env)
		}
	}

	if len(missingEnvs) > 0 {
		log.Fatalf("Required environment variables are missing: %v", missingEnvs)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for environment variable %s, using default %d. Error: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("Warning: Invalid value for environment variable %s, using default %d. Error: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for environment variable %s, using default %v. Error: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid value for environment variable %s, using default %v. Error: %v", key, defaultValue, err)
		return defaultValue
	}
	return value
}

func Get() *Config {
	if GlobalConfig == nil {
		log.Fatal("Config not initialized. Call InitConfig() first.")
	}
	return GlobalConfig
}
