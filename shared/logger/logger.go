// Package logger provides configurable logging for both application and Fiber HTTP requests
// with support for different log levels using uber-go/zap logger.
//
// Usage:
//
//	if err := logger.Init(); err != nil {
//	    log.Fatal("Failed to initialize logger:", err)
//	}
//	defer logger.Sync()
//
//	logger.Info("Application started")
package logger

import (
	"fmt"
	"go-api/config"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Global logger instances
var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Close gracefully closes the logger
func Close() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

const (
	defaultMaxSize    = 10 // 10MB
	defaultMaxBackups = 5
)

// Init initializes the logger with file output using zap
func Init() error {
	cfg := config.Get()

	// Create log directory if it doesn't exist (hardcoded to storage/logs)
	logDir := "storage/logs"
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Setup log file with rotation (single file only)
	logPath := filepath.Join(logDir, "app.log")

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    defaultMaxSize, // Use hardcoded 10MB default
		MaxBackups: defaultMaxBackups,
		Compress:   true,
		// MaxAge removed - files will never be deleted based on age
	}

	// Create cores for console and file output
	consoleCore := createConsoleCore(cfg.LogLevel)
	fileCore := createFileCore(lumberjackLogger, cfg.LogLevel)

	// Combine cores
	core := zapcore.NewTee(consoleCore, fileCore)

	// Create logger with caller info and stack trace for errors
	Logger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)
	
	if Logger == nil {
		return fmt.Errorf("failed to initialize zap logger")
	}
	
	Sugar = Logger.Sugar()
	
	if Sugar == nil {
		return fmt.Errorf("failed to initialize zap sugar logger")
	}

	return nil
}

// createConsoleCore creates a console output core with colored output
func createConsoleCore(logLevel string) zapcore.Core {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(os.Stdout),
		parseLogLevel(logLevel),
	)
}

// createFileCore creates a file output core with JSON format
func createFileCore(writer io.Writer, logLevel string) zapcore.Core {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.LevelKey = "level"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.CallerKey = "caller"
	config.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(writer),
		parseLogLevel(logLevel),
	)
}

// parseLogLevel converts string log level to zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// GetAppLoggerConfig returns application logger middleware configuration
func GetAppLoggerConfig() logger.Config {
	cfg := config.Get()

	var output io.Writer = os.Stdout

	// Setup file output if enabled with proper error handling
	if cfg.EnableAppLog {
		logDir := "storage/logs"
		if err := os.MkdirAll(logDir, 0750); err != nil {
			// Fallback to stdout only if directory creation fails
			output = os.Stdout
		} else {
			// Use hardcoded filename for security
			const logFileName = "app-access.log"
			logPath := filepath.Join(logDir, logFileName)
			
			// Additional safety check to ensure we're within the expected directory
			if !strings.HasPrefix(logPath, logDir) {
				output = os.Stdout
			} else {
				// #nosec G304 -- path is validated and controlled
				if logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
					// Fallback to stdout only if file opening fails
					output = os.Stdout
				} else {
					// Note: We don't close the file here as it's used by the middleware
					// The file will be closed when the application shuts down
					output = io.MultiWriter(os.Stdout, logFile)
				}
			}
		}
	}

	return logger.Config{
		Output:     output,
		Format:     getAppLogFormat(cfg.LogLevel),
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Next:       getAppLogFilter(cfg.LogLevel),
	}
}

// getAppLogFormat returns log format based on log level
func getAppLogFormat(logLevel string) string {
	formats := map[string]string{
		"debug": "${time} [${level}] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${query} | ${ua} | ${error}\n",
		"info":  "${time} [INFO] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		"warn":  "${time} [${level}] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		"error": "${time} [ERROR] | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}

	if format, exists := formats[logLevel]; exists {
		return format
	}
	return formats["info"] // default
}

// getAppLogFilter returns filter function based on log level
func getAppLogFilter(logLevel string) func(*fiber.Ctx) bool {
	staticExtensions := []string{".css", ".js", ".ico", ".png", ".jpg", ".svg"}

	switch logLevel {
	case "debug":
		return nil // Log everything
	case "info":
		return func(c *fiber.Ctx) bool {
			// Skip successful static file requests
			if c.Response().StatusCode() == 200 {
				for _, ext := range staticExtensions {
					if filepath.Ext(c.Path()) == ext {
						return true
					}
				}
			}
			return false
		}
	case "warn":
		return func(c *fiber.Ctx) bool {
			return c.Response().StatusCode() < 400 // Only log 4xx and 5xx
		}
	case "error":
		return func(c *fiber.Ctx) bool {
			return c.Response().StatusCode() < 500 // Only log 5xx
		}
	default:
		return getAppLogFilter("info")
	}
}

// Logging functions with zap fields
func Info(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
	}
}

func Error(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(message, fields...)
	}
}

func Debug(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(message, fields...)
	}
}

func Warn(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(message, fields...)
	}
}

func Fatal(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(message, fields...)
	}
}

// Sugar logger convenience functions
func Infof(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Infof(template, args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Errorf(template, args...)
	}
}

func Debugf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Debugf(template, args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Warnf(template, args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Fatalf(template, args...)
	}
}

// Field helper functions for structured logging
func String(key, val string) zap.Field                 { return zap.String(key, val) }
func Int(key string, val int) zap.Field                { return zap.Int(key, val) }
func Int64(key string, val int64) zap.Field            { return zap.Int64(key, val) }
func Float64(key string, val float64) zap.Field        { return zap.Float64(key, val) }
func Bool(key string, val bool) zap.Field              { return zap.Bool(key, val) }
func Duration(key string, val time.Duration) zap.Field { return zap.Duration(key, val) }
func Time(key string, val time.Time) zap.Field         { return zap.Time(key, val) }
func Err(err error) zap.Field                          { return zap.Error(err) }
func Any(key string, val interface{}) zap.Field        { return zap.Any(key, val) }

// With adds fields to the logger and returns a new logger
func With(fields ...zap.Field) *zap.Logger {
	if Logger != nil {
		return Logger.With(fields...)
	}
	return nil
}

// Named creates a named logger
func Named(name string) *zap.Logger {
	if Logger != nil {
		return Logger.Named(name)
	}
	return nil
}

// safePath validates and sanitizes file paths to prevent directory traversal
func safePath(basePath, filename string) (string, error) {
	// Clean the filename to remove any path traversal attempts
	cleanFilename := filepath.Clean(filename)
	
	// Check for path traversal attempts
	if strings.Contains(cleanFilename, "..") {
		return "", fmt.Errorf("invalid filename: path traversal detected")
	}
	
	// Ensure filename doesn't contain absolute path indicators
	if filepath.IsAbs(cleanFilename) {
		return "", fmt.Errorf("invalid filename: absolute paths not allowed")
	}
	
	// Join with base path and clean again
	fullPath := filepath.Join(basePath, cleanFilename)
	
	// Ensure the resulting path is still within the base directory
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base path: %w", err)
	}
	
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve full path: %w", err)
	}
	
	if !strings.HasPrefix(absPath, absBase) {
		return "", fmt.Errorf("path outside of allowed directory")
	}
	
	return fullPath, nil
}
