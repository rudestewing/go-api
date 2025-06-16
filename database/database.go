package database

import (
	"fmt"
	"go-api/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDatabaseLogger returns a configured database logger based on config settings
func getDatabaseLogger() logger.Interface {
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


// InitDB initializes the database connection with GORM
func InitDB() (*gorm.DB, error) {
	cfg := config.Get()

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// Configure database logger based on config
	databaseLogger := getDatabaseLogger()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: databaseLogger,
	})

	if err != nil {
		return nil, err
	}

	// Configure connection pool to prevent memory leaks
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings from config
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)

	// Auto migrate the user model
	// if err := db.AutoMigrate(&model.User{}); err != nil {
	// 	return nil, fmt.Errorf("failed to migrate database: %w", err)
	// }

	log.Println("Database connection established")
	return db, nil
}
