package database

import (
	"go-api/config"

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

	// Use the default logger with proper configuration
	return logger.Default.LogMode(logLevel)
}
