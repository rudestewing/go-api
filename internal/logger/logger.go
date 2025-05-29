package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

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
func InitLogger(config LoggerConfig) error {
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

// GetFiberLoggerConfig returns Fiber logger middleware configuration
func GetFiberLoggerConfig(config LoggerConfig) logger.Config {
	// Create log file for Fiber specifically
	logFileName := "fiber-access.log"
	if config.EnableDaily {
		logFileName = fmt.Sprintf("fiber-access-%s.log", time.Now().Format("2006-01-02"))
	}

	logPath := filepath.Join(config.LogDir, logFileName)

	// Open or create Fiber log file
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open Fiber log file: %v", err)
		// Fallback to stdout only
		return logger.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		}
	}

	// Create multi-writer for Fiber logs (file + stdout)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return logger.Config{
		Output: multiWriter,
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone: "Local",
	}
}

// LogInfo logs info level messages
func LogInfo(message string, args ...any) {
	log.Printf("[INFO] "+message, args...)
}

// LogError logs error level messages
func LogError(message string, args ...any) {
	log.Printf("[ERROR] "+message, args...)
}

// LogWarning logs warning level messages
func LogWarning(message string, args ...any) {
	log.Printf("[WARNING] "+message, args...)
}

// LogDebug logs debug level messages
func LogDebug(message string, args ...any) {
	log.Printf("[DEBUG] "+message, args...)
}

// LogFatal logs fatal level messages and exits
func LogFatal(message string, args ...any) {
	log.Fatalf("[FATAL] "+message, args...)
}

// CleanupOldLogs removes log files older than maxAge days
func CleanupOldLogs(config LoggerConfig) error {
	if !config.EnableDaily {
		return nil
	}

	return filepath.Walk(config.LogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file is older than maxAge
		if time.Since(info.ModTime()).Hours() > float64(config.MaxAge*24) {
			log.Printf("Removing old log file: %s", path)
			return os.Remove(path)
		}

		return nil
	})
}
