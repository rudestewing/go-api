package main

import (
	"context"
	"go-api/config"
	"go-api/container"
	"go-api/handler"
	"go-api/logger"
	"go-api/router"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func createFiberApp() *fiber.App {
	cfg := config.Get()

	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Log error to file
			logger.LogError("Fiber error: %s | Path: %s | Method: %s | IP: %s",
				err.Error(), c.Path(), c.Method(), c.IP())

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())

	// Setup logger configuration from config
	loggerConfig := logger.LoggerConfig{
		LogDir:      cfg.LogDir,
		MaxSize:     cfg.LogMaxSize,
		MaxAge:      cfg.LogMaxAge,
		EnableDaily: cfg.EnableDailyLog,
	}

	// Use custom Fiber logger with file output
	fiberLoggerConfig := logger.GetFiberLoggerConfig(loggerConfig)
	app.Use(fiberLogger.New(fiberLoggerConfig))

	app.Use(cors.New())

	return app
}

func main() {
	// Initialize config first
	config.InitConfig()

	// Initialize logger
	cfg := config.Get()
	loggerConfig := logger.LoggerConfig{
		LogDir:      cfg.LogDir,
		MaxSize:     cfg.LogMaxSize,
		MaxAge:      cfg.LogMaxAge,
		EnableDaily: cfg.EnableDailyLog,
	}

	if err := logger.InitLogger(loggerConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Clean up old logs on startup
	if err := logger.CleanupOldLogs(loggerConfig); err != nil {
		logger.LogWarning("Failed to cleanup old logs: %v", err)
	}

	logger.LogInfo("Starting application...")

	// Initialize Fiber App
	app := createFiberApp()

	// Initialize App Container
	container, err_container := container.NewContainer()
	if err_container != nil {
		logger.LogFatal("Failed to create container: %v", err_container)
	}

	// Initialize Handler
	handler := handler.NewHandler(container)

	// Setup Routes
	router.RegisterRoutes(app, handler)

	port := ":" + cfg.AppPort
	logger.LogInfo("🚀 Server starting on port %s...", cfg.AppPort)

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(port); err != nil {
			logger.LogFatal("Server error: %v", err)
		}
	}()

	logger.LogInfo("Server started successfully. Press Ctrl+C to shutdown gracefully...")

	// Block until we receive an interrupt signal
	<-c

	logger.LogInfo("Shutting down server gracefully...")

	// Create a deadline for the shutdown using config
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	// Shutdown the Fiber server
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.LogError("Server forced to shutdown: %v", err)
	}

	// Close container (database connections, etc.)
	if err := container.Close(ctx); err != nil {
		logger.LogError("Error closing container: %v", err)
	}

	logger.LogInfo("Server exited successfully")
}
