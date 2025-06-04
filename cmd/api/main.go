package main

import (
	"context"
	"go-api/app/handler"
	"go-api/config"
	"go-api/container"
	"go-api/middleware"
	"go-api/router"
	"go-api/shared/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func createFiberApp() *fiber.App {
	cfg := config.Get()

	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		IdleTimeout:           cfg.IdleTimeout,
		DisableStartupMessage: !cfg.EnableFiberLog, // Disable startup message jika fiber log disabled
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			// Log error to file
			logger.Errorf("Fiber error: %s | Path: %s | Method: %s | IP: %s",
				err.Error(), c.Path(), c.Method(), c.IP())

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware order is important!
	app.Use(recover.New())

	// Request ID for tracing
	app.Use(middleware.RequestIDMiddleware())

	// Error handling (should be early in the chain)
	app.Use(middleware.ErrorHandlingMiddleware())

	// Configuration validation
	app.Use(middleware.ConfigValidationMiddleware())

	// Input validation middleware (before other processing)
	app.Use(middleware.ContentTypeValidationMiddleware())
	app.Use(middleware.InputValidationMiddleware())

	// Security Headers (helmet)
	if cfg.SecurityHeadersEnabled {
		app.Use(helmet.New())
	}

	// Rate Limiting (general)
	if cfg.RateLimitEnabled {
		app.Use(middleware.RateLimitMiddleware())
	}

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})
	// Use custom Fiber logger with file output - hanya jika enabled
	if cfg.EnableFiberLog {
		fiberLoggerConfig := logger.GetFiberConfig()
		app.Use(fiberLogger.New(fiberLoggerConfig))
	}

	// CORS with production-safe configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		AllowCredentials: false,
	}))

	return app
}

func main() {
	// Initialize config first
	config.InitConfig()

	cfg := config.Get()
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Clean up old logs on startup
	if err := logger.CleanupOldLogs(); err != nil {
		logger.Warnf("Failed to cleanup old logs: %v", err)
	}

	logger.Infof("Starting application in %s environment...", cfg.Environment)

	// Initialize Fiber App
	app := createFiberApp()

	// Initialize App Container
	container, err_container := container.NewContainer()

	if err_container != nil {
		logger.Fatalf("Failed to create container: %v", err_container)
	}

	handler := handler.NewHandler(container)
	// Setup Routes
	router.RegisterRoutes(app, handler, container)

	port := ":" + cfg.AppPort
	logger.Infof("🚀 Server starting on port %s...", cfg.AppPort)

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine with error handling
	serverErr := make(chan error, 1)
	go func() {
		if err := app.Listen(port); err != nil {
			logger.Errorf("Server failed to start: %v", err)
			serverErr <- err
		}
	}()

	logger.Infof("Server started successfully. Press Ctrl+C to shutdown gracefully...")

	// Block until we receive an interrupt signal or server error
	select {
	case <-c:
		logger.Infof("Received interrupt signal")
	case err := <-serverErr:
		logger.Errorf("Server error: %v", err)
	}

	logger.Infof("Shutting down server gracefully...")

	// Create a deadline for the shutdown using config with fallback
	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second // Default fallback
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown the Fiber server
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}
	// Close container (database connections, etc.)
	if err := container.Close(ctx); err != nil {
		logger.Errorf("Error closing container: %v", err)
	}

	logger.Infof("Server exited successfully")
}
