package database

import (
	"context"
	"fmt"
	"go-api/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDatabaseLogger returns a configured database logger based on config settings
func GetDatabaseLogger() logger.Interface {
	cfg := config.Get()

	if !cfg.EnableDatabaseLog {
		return logger.Default.LogMode(logger.Silent)
	}

	// Determine log level based on config
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "debug":
		logLevel = logger.Info // Database's most verbose level
	case "info":
		logLevel = logger.Warn // Only slow queries and errors
	case "warn":
		logLevel = logger.Error // Only errors
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Warn
	}

	// Use the default logger with proper configuration
	return logger.Default.LogMode(logLevel)
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	return initDatabaseConnection(cfg)
}

// initDatabaseConnection is the internal function that handles the actual database connection
func initDatabaseConnection(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL configuration is required")
	}

	// Configure database logger based on config
	databaseLogger := GetDatabaseLogger()

	// Add connection timeout and retry logic
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: databaseLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Enable prepared statements for better performance
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool to prevent memory leaks and optimize performance
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings from config with validation
	if cfg.DBMaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	}
	if cfg.DBMaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// InitDB initializes a database connection using the global config
// This function is used by seeders and other standalone utilities
func InitDB() (*gorm.DB, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("configuration not initialized")
	}
	
	return NewConnection(cfg)
}
