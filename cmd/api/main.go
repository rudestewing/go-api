package main

import (
	"context"
	"go-api/app/router"
	"go-api/app/shared/logger"
	"go-api/config"
	"go-api/container"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
			logger.LogError("Fiber error: %s | Path: %s | Method: %s | IP: %s",
				err.Error(), c.Path(), c.Method(), c.IP())

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())

	// Security Headers (helmet)
	if cfg.SecurityHeadersEnabled {
		app.Use(helmet.New())
	}

	// Rate Limiting
	if cfg.RateLimitEnabled {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimitMax,
			Expiration: cfg.RateLimitWindow,
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(429).JSON(fiber.Map{
					"error": "Too many requests, please try again later",
				})
			},
		}))
	}

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})

	// Setup logger configuration from config
	loggerConfig := logger.LoggerConfig{
		LogDir:      cfg.LogDir,
		MaxSize:     cfg.LogMaxSize,
		MaxAge:      cfg.LogMaxAge,
		EnableDaily: cfg.EnableDailyLog,
	}

	// Use custom Fiber logger with file output - hanya jika enabled
	if cfg.EnableFiberLog {
		fiberLoggerConfig := logger.GetFiberLoggerConfig(loggerConfig)
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

	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Clean up old logs on startup
	if err := logger.CleanupOldLogs(); err != nil {
		logger.LogWarning("Failed to cleanup old logs: %v", err)
	}

	logger.LogInfo("Starting application in %s environment...", cfg.Environment)

	// Initialize Fiber App
	app := createFiberApp()

	// Initialize App Container
	container, err_container := container.NewContainer()

	if err_container != nil {
		logger.LogFatal("Failed to create container: %v", err_container)
	}

	// Setup Routes
	router.RegisterRoutes(app, container)

	port := ":" + cfg.AppPort
	logger.LogInfo("🚀 Server starting on port %s...", cfg.AppPort)

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine with error handling
	serverErr := make(chan error, 1)
	go func() {
		if err := app.Listen(port); err != nil {
			logger.LogError("Server failed to start: %v", err)
			serverErr <- err
		}
	}()

	logger.LogInfo("Server started successfully. Press Ctrl+C to shutdown gracefully...")

	// Block until we receive an interrupt signal or server error
	select {
	case <-c:
		logger.LogInfo("Received interrupt signal")
	case err := <-serverErr:
		logger.LogError("Server error: %v", err)
	}

	logger.LogInfo("Shutting down server gracefully...")

	// Create a deadline for the shutdown using config with fallback
	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second // Default fallback
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
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
