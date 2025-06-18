package cmd

import (
	"context"
	"go-api/app"
	"go-api/config"
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
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long: `Start the HTTP API server with Fiber framework.

This command will:
- Initialize configuration and logger
- Start the Fiber web server
- Setup middleware and routes
- Handle graceful shutdown on interrupt signals

Examples:
  serve`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}


func init() {
	RootCmd.AddCommand(serveCmd)
}


func createFiberApp(provider *app.Provider) *fiber.App {
	cfg := config.Get()

	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		IdleTimeout:           cfg.IdleTimeout,
		DisableStartupMessage: !cfg.EnableAppLog, // Disable startup message if app log disabled
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
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Environment == "development",
	}))

	// Request ID for tracing
	app.Use(middleware.RequestIDMiddleware())

	// Security middleware
	app.Use(middleware.SecurityHeadersMiddleware())
	app.Use(middleware.InputSanitizationMiddleware())
	app.Use(middleware.RequestValidationMiddleware())

	// Timeout middleware
	app.Use(middleware.TimeoutMiddleware(cfg.ReadTimeout, "Request timeout"))

	// Error handling (should be early in the chain)
	app.Use(middleware.ErrorHandlingMiddleware())

	// Configuration validation
	app.Use(middleware.ConfigValidationMiddleware())

	// Input validation middleware (before other processing)
	app.Use(middleware.ContentTypeValidationMiddleware())
	app.Use(middleware.InputValidationMiddleware())

	// Security Headers (helmet) - additional security
	if cfg.SecurityHeadersEnabled {
		app.Use(helmet.New(helmet.Config{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "DENY",
			ReferrerPolicy:        "no-referrer",
			CrossOriginEmbedderPolicy: "require-corp",
		}))
	}

	// Rate Limiting (general)
	if cfg.RateLimitEnabled {
		app.Use(middleware.AdvancedRateLimitMiddleware())
	}

	// Health check endpoints
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check actual database status via AppService
		dbStatus := "connected"
		if provider.DB != nil {
			if sqlDB, err := provider.DB.DB(); err != nil || sqlDB.Ping() != nil {
				dbStatus = "disconnected"
			}
		} else {
			dbStatus = "not_initialized"
		}
		
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
			"service":   "go-api",
			"database":  dbStatus,
		})
	})
	
	app.Get("/ready", func(c *fiber.Ctx) error {
		// Quick readiness check without database
		return c.JSON(fiber.Map{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})
	
	// Use custom application logger with file output - only if enabled
	if cfg.EnableAppLog {
		appLoggerConfig := logger.GetAppLoggerConfig()
		app.Use(fiberLogger.New(appLoggerConfig))
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

func startServer() {
	// Initialize config first
	config.InitConfig()

	cfg := config.Get()

	// Initialize logger with proper error handling
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	defer logger.Sync() // Ensure logs are flushed on exit
	
	logger.Infof("Starting application in %s environment...", cfg.Environment)
	// Initialize global database first
	logger.Infof("Initializing application services...")
	provider, err := app.BootProvider(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize application services: %v", err)
	}
	logger.Infof("✅ Application services initialized successfully")

	// Initialize Fiber App
	fiberApp := createFiberApp(provider)

	
	// Setup Routes
	router.RegisterRoutes(fiberApp, provider)

	port := ":" + cfg.AppPort
	logger.Infof("🚀 Server starting on port %s...", cfg.AppPort)

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine with error handling
	serverErr := make(chan error, 1)
	go func() {
		if err := fiberApp.Listen(port); err != nil {
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

	// Shutdown the Fiber server gracefully
	logger.Infof("Shutting down server...")
	
	if err := fiberApp.ShutdownWithContext(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	} else {
		logger.Infof("Server shutdown completed")
	}

	// Close global database connection
	provider.ShutdownProvider()

	// Final log sync
	logger.Sync()
	logger.Infof("Application exited successfully")
}
