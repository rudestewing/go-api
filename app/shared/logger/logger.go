// Package logger provides configurable logging for both application and Fiber HTTP requests
// with support for different log levels similar to database logging.
//
// Log Levels:
// - debug: Most verbose, logs everything including detailed request info
// - info:  Standard logging with essential information
// - warn:  Only warnings (4xx) and errors (5xx)
// - error: Only server errors (5xx)
//
// Usage:
//
//	cfg := config.Get()
//	if cfg.EnableFiberLog {
//	    app.Use(fiberLogger.New(logger.GetFiberLoggerConfig(loggerConfig)))
//	}
package logger

import (
	"fmt"
	"go-api/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// LoggerConfig holds the configuration for logging
type LoggerConfig struct {
	LogDir      string
	MaxSize     int64 // in bytes
	MaxAge      int   // in days
	EnableDaily bool  // create daily log files
}

// DefaultConfig returns default logger configuration
func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		LogDir:      "storage/logs",
		MaxSize:     10 * 1024 * 1024, // 10MB
		MaxAge:      30,                // 30 days
		EnableDaily: true,
	}
}

// InitLogger initializes the logger with file output
func InitLogger() error {
	cfg := config.Get()

	config := LoggerConfig{
		LogDir:      cfg.LogDir,
		EnableDaily: cfg.EnableDailyLog,
	}
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with current date
	logFileName := "app.log"
	if config.EnableDaily {
		logFileName = fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	}

	logPath := filepath.Join(config.LogDir, logFileName)

	// Open or create log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Set up multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return nil
}

// GetFiberLoggerConfig returns Fiber logger middleware configuration with log level support
func GetFiberLoggerConfig(loggerConfig LoggerConfig) logger.Config {
	cfg := config.Get()

	// Create log file for Fiber specifically
	logFileName := "fiber-access.log"
	if loggerConfig.EnableDaily {
		logFileName = fmt.Sprintf("fiber-access-%s.log", time.Now().Format("2006-01-02"))
	}

	logPath := filepath.Join(loggerConfig.LogDir, logFileName)

	// Open or create Fiber log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open Fiber log file: %v", err)
		// Fallback to stdout only
		return getBasicFiberLoggerConfig(cfg.LogLevel)
	}

	// Create multi-writer for Fiber logs (file + stdout)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return logger.Config{
		Output:     multiWriter,
		Format:     getFiberLogFormat(cfg.LogLevel),
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Next:       getFiberLogFilter(cfg.LogLevel),
	}
}

// GetFiberLoggerWithLevel returns a simple Fiber logger config for a specific log level
// This is a convenience function for quick setup without file logging
func GetFiberLoggerWithLevel(logLevel string) logger.Config {
	return logger.Config{
		Format:     getFiberLogFormat(logLevel),
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Next:       getFiberLogFilter(logLevel),
	}
}

// getBasicFiberLoggerConfig returns basic logger config without file output
func getBasicFiberLoggerConfig(logLevel string) logger.Config {
	return logger.Config{
		Format:     getFiberLogFormat(logLevel),
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Next:       getFiberLogFilter(logLevel),
	}
}

// getFiberLogFormat returns log format based on log level
func getFiberLogFormat(logLevel string) string {
	switch logLevel {
	case "debug":
		// Most verbose - include all fields including user agent and headers
		return "${time} [${level}] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${query} | ${ua} | ${error}\n"
	case "info":
		// Standard format - essential information
		return "${time} [INFO] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n"
	case "warn":
		// Only warnings and errors (4xx, 5xx)
		return "${time} [${level}] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n"
	case "error":
		// Only errors (5xx)
		return "${time} [ERROR] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n"
	default:
		// Default to info level
		return "${time} [INFO] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n"
	}
}

// getFiberLogFilter returns filter function based on log level
func getFiberLogFilter(logLevel string) func(*fiber.Ctx) bool {
	switch logLevel {
	case "debug":
		// Log everything
		return nil
	case "info":
		// Skip only successful static file requests to reduce noise
		return func(c *fiber.Ctx) bool {
			// Skip logging for static files with 200 status
			if c.Response().StatusCode() == 200 &&
			   (filepath.Ext(c.Path()) == ".css" ||
			    filepath.Ext(c.Path()) == ".js" ||
			    filepath.Ext(c.Path()) == ".ico" ||
			    filepath.Ext(c.Path()) == ".png" ||
			    filepath.Ext(c.Path()) == ".jpg" ||
			    filepath.Ext(c.Path()) == ".svg") {
				return true
			}
			return false
		}
	case "warn":
		// Only log warnings (4xx) and errors (5xx)
		return func(c *fiber.Ctx) bool {
			status := c.Response().StatusCode()
			return status < 400
		}
	case "error":
		// Only log errors (5xx)
		return func(c *fiber.Ctx) bool {
			status := c.Response().StatusCode()
			return status < 500
		}
	default:
		// Default to info level behavior
		return func(c *fiber.Ctx) bool {
			if c.Response().StatusCode() == 200 &&
			   (filepath.Ext(c.Path()) == ".css" ||
			    filepath.Ext(c.Path()) == ".js" ||
			    filepath.Ext(c.Path()) == ".ico" ||
			    filepath.Ext(c.Path()) == ".png" ||
			    filepath.Ext(c.Path()) == ".jpg" ||
			    filepath.Ext(c.Path()) == ".svg") {
				return true
			}
			return false
		}
	}
}

// LogInfo logs info level messages - respects log level setting
func LogInfo(message string, args ...any) {
	cfg := config.Get()
	if shouldLog("info", cfg.LogLevel) {
		log.Printf("[INFO] "+message, args...)
	}
}

// LogError logs error level messages - always logged
func LogError(message string, args ...any) {
	log.Printf("[ERROR] "+message, args...)
}

// LogWarning logs warning level messages - respects log level setting
func LogWarning(message string, args ...any) {
	cfg := config.Get()
	if shouldLog("warn", cfg.LogLevel) {
		log.Printf("[WARNING] "+message, args...)
	}
}

// LogDebug logs debug level messages - respects log level setting
func LogDebug(message string, args ...any) {
	cfg := config.Get()
	if shouldLog("debug", cfg.LogLevel) {
		log.Printf("[DEBUG] "+message, args...)
	}
}

// LogFatal logs fatal level messages and exits - always logged
func LogFatal(message string, args ...any) {
	log.Fatalf("[FATAL] "+message, args...)
}

// shouldLog determines if a message should be logged based on log level
func shouldLog(messageLevel, configLevel string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	msgLevel, exists := levels[messageLevel]
	if !exists {
		return true // Log unknown levels
	}

	cfgLevel, exists := levels[configLevel]
	if !exists {
		cfgLevel = 1 // Default to info level
	}

	return msgLevel >= cfgLevel
}

// CleanupOldLogs removes log files older than maxAge days
func CleanupOldLogs() error {
	cfg := config.Get()

	loggerConfig := LoggerConfig{
		LogDir:      cfg.LogDir,
		MaxAge:      cfg.LogMaxAge,
		EnableDaily: cfg.EnableDailyLog,
	}

	if !loggerConfig.EnableDaily {
		return nil
	}

	return filepath.Walk(loggerConfig.LogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is older than maxAge
		if time.Since(info.ModTime()).Hours() > float64(loggerConfig.MaxAge*24) {
			log.Printf("Removing old log file: %s", path)
			return os.Remove(path)
		}

		return nil
	})
}
