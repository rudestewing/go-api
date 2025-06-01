package database

import (
	"go-api/config"
	"time"

	"gorm.io/gorm/logger"
)

// GetGormLogger returns a configured GORM logger based on config settings
func GetGormLogger() logger.Interface {
	cfg := config.Get()

	if !cfg.EnableGormLog {
		return logger.Default.LogMode(logger.Silent)
	}

	// Determine log level based on config
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "debug":
		logLevel = logger.Info // GORM's most verbose level
	case "info":
		logLevel = logger.Warn // Only slow queries and errors
	case "warn":
		logLevel = logger.Error // Only errors
	case "error":
		logLevel = logger.Error
	default:
		logLevel = logger.Warn
	}

	// Custom logger configuration
	return logger.New(
		nil, // Use default writer (stdout)
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Log slow queries
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,  // Don't log "record not found" errors
			Colorful:                  true,  // Enable colorful output
			ParameterizedQueries:      false, // Log full SQL with parameters
		},
	)
}
