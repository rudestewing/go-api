// Package logger provides configurable logging for both application and Fiber HTTP requests
// with support for different log levels using uber-go/zap logger.
//
// Log Levels:
// - debug: Most verbose, logs everything including detailed request info
// - info:  Standard logging with essential information
// - warn:  Only warnings (4xx) and errors (5xx)
// - error: Only server errors (5xx)
//
// Usage:
//
//	if err := logger.InitLogger(); err != nil {
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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Global zap logger instance
var Logger *zap.Logger
var Sugar *zap.SugaredLogger

// LoggerConfig holds the configuration for logging
type LoggerConfig struct {
	LogDir      string
	MaxSize     int   // in megabytes
	MaxAge      int   // in days
	MaxBackups  int   // number of backup files
	EnableDaily bool  // create daily log files
	Compress    bool  // compress old log files
}

// DefaultConfig returns default logger configuration
func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		LogDir:      "storage/logs",
		MaxSize:     10, // 10MB (lumberjack uses MB units)
		MaxAge:      30, // 30 days
		MaxBackups:  5,
		EnableDaily: true,
		Compress:    true,
	}
}

// InitLogger initializes the logger with file output using zap
func InitLogger() error {
	cfg := config.Get()

	loggerConfig := LoggerConfig{
		LogDir:      cfg.LogDir,
		MaxSize:     10, // 10MB
		MaxAge:      cfg.LogMaxAge,
		MaxBackups:  5,
		EnableDaily: cfg.EnableDailyLog,
		Compress:    true,
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(loggerConfig.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Setup log rotation with lumberjack
	logFileName := "app.log"
	if loggerConfig.EnableDaily {
		logFileName = fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	}

	logPath := filepath.Join(loggerConfig.LogDir, logFileName)

	// Configure lumberjack for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    loggerConfig.MaxSize,
		MaxAge:     loggerConfig.MaxAge,
		MaxBackups: loggerConfig.MaxBackups,
		Compress:   loggerConfig.Compress,
	}
	// Create writeSyncer that writes to both file and stdout
	fileWriter := zapcore.AddSync(lumberjackLogger)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Configure log level based on config
	logLevel := parseLogLevel(cfg.LogLevel)

	// Create encoder configs for console and file
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.TimeKey = "timestamp"
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoderConfig.LevelKey = "level"
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoderConfig.CallerKey = "caller"
	fileEncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Create separate cores for console and file
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		consoleWriter,
		logLevel,
	)

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoderConfig),
		fileWriter,
		logLevel,
	)

	// Combine cores
	core := zapcore.NewTee(consoleCore, fileCore)

	// Create logger with caller info and stack trace for errors
	Logger = zap.New(core, 
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1), // Skip one level for wrapper functions
	)
	Sugar = Logger.Sugar()

	return nil
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
		Sugar.Errorf("Failed to open Fiber log file: %v", err)
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

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// Zap field wrapper functions for structured logging
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Error(key string, err error) zap.Field {
	return zap.Error(err)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// Additional field helper functions for common use cases
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}

// Info logs info level messages using zap
func Info(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(message, fields...)
	}
}

// Error logs error level messages using zap
func ErrorLog(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(message, fields...)
	}
}

// Debug logs debug level messages using zap
func Debug(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(message, fields...)
	}
}

// Warn logs warning level messages using zap
func Warn(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(message, fields...)
	}
}

// Fatal logs fatal level messages and exits using zap
func Fatal(message string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(message, fields...)
	}
}

// Sugar logger convenience functions
func Infof(message string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Infof(message, args...)
	}
}

func Errorf(message string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Errorf(message, args...)
	}
}

func Debugf(message string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Debugf(message, args...)
	}
}

func Warnf(message string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Warnf(message, args...)
	}
}

func Fatalf(message string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Fatalf(message, args...)
	}
}

// With adds fields to the logger and returns a new logger
func With(fields ...zap.Field) *zap.Logger {
	if Logger != nil {
		return Logger.With(fields...)
	}
	return nil
}

// WithContext returns a logger with context fields
func WithContext(ctx map[string]interface{}) *zap.SugaredLogger {
	if Sugar != nil {
		var fields []interface{}
		for k, v := range ctx {
			fields = append(fields, k, v)
		}
		return Sugar.With(fields...)
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

// NamedSugar creates a named sugar logger
func NamedSugar(name string) *zap.SugaredLogger {
	if Sugar != nil {
		return Sugar.Named(name)
	}
	return nil
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
			if Sugar != nil {
				Sugar.Infof("Removing old log file: %s", path)
			}
			return os.Remove(path)
		}

		return nil
	})
}
