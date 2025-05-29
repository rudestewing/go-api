package main

import (
	"go-api/config"
	"go-api/container"
	"go-api/internal/handler"
	"go-api/internal/logger"
	"go-api/router"
	"log"

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

	// Serve App
	if err := app.Listen(port); err != nil {
		logger.LogFatal("Server error: %v", err)
	}
}
