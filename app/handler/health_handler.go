package handler

import (
	"go-api/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck provides a simple health check endpoint
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": config.Now(),
		"service":   "go-api",
	})
}

// DetailedHealthCheck provides detailed health information
func (h *HealthHandler) DetailedHealthCheck(c *fiber.Ctx) error {
	health := fiber.Map{
		"status":    "healthy",
		"timestamp": config.Now(),
		"service":   "go-api",
		"version":   "1.0.0", // You can make this configurable
		"checks":    fiber.Map{},
	}

	checks := health["checks"].(fiber.Map)

	// Database health check
	dbHealth := h.checkDatabase()
	checks["database"] = dbHealth

	// Overall status based on checks
	overallHealthy := true
	for _, check := range checks {
		if checkMap, ok := check.(fiber.Map); ok {
			if status, exists := checkMap["status"]; exists && status != "healthy" {
				overallHealthy = false
				break
			}
		}
	}

	if !overallHealthy {
		health["status"] = "unhealthy"
		return c.Status(fiber.StatusServiceUnavailable).JSON(health)
	}

	return c.Status(fiber.StatusOK).JSON(health)
}

// ReadinessCheck checks if the service is ready to accept requests
func (h *HealthHandler) ReadinessCheck(c *fiber.Ctx) error {
	// Check database connectivity
	dbCheck := h.checkDatabase()
	if dbCheck["status"] != "healthy" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "not_ready",
			"message": "Database is not available",
			"details": dbCheck,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "ready",
		"timestamp": config.Now(),
	})
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "alive",
		"timestamp": config.Now(),
	})
}

// checkDatabase verifies database connectivity and basic operations
func (h *HealthHandler) checkDatabase() fiber.Map {
	start := time.Now()

	// Get the underlying SQL database
	sqlDB, err := h.db.DB()
	if err != nil {
		return fiber.Map{
			"status":       "unhealthy",
			"error":        "Failed to get database instance",
			"response_time": time.Since(start).String(),
		}
	}

	// Test database connectivity
	if err := sqlDB.Ping(); err != nil {
		return fiber.Map{
			"status":       "unhealthy",
			"error":        "Database ping failed",
			"response_time": time.Since(start).String(),
		}
	}

	// Test basic query
	var result int
	if err := h.db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fiber.Map{
			"status":       "unhealthy",
			"error":        "Database query failed",
			"response_time": time.Since(start).String(),
		}
	}

	// Get database stats
	stats := sqlDB.Stats()

	return fiber.Map{
		"status":        "healthy",
		"response_time": time.Since(start).String(),
		"connections": fiber.Map{
			"open":        stats.OpenConnections,
			"in_use":      stats.InUse,
			"idle":        stats.Idle,
			"max_open":    stats.MaxOpenConnections,
		},
	}
}
