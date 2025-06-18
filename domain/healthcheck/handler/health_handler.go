package handler

import (
	"context"
	"database/sql"
	"go-api/app"
	"go-api/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(p *app.Provider) *HealthHandler {
	return &HealthHandler{
		db: p.DB,
	}
}

// HealthCheck provides basic health status
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"service":   "go-api",
	})
}

// ReadinessCheck checks if the service is ready to serve requests
func (h *HealthHandler) ReadinessCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	
	checks := make(map[string]interface{})
	overallStatus := "ok"
	
	// Database connectivity check
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			checks["database"] = map[string]interface{}{
				"status": "error",
				"error":  "Failed to get database instance",
			}
			overallStatus = "error"
		} else {
			if err := sqlDB.PingContext(ctx); err != nil {
				checks["database"] = map[string]interface{}{
					"status": "error",
					"error":  err.Error(),
				}
				overallStatus = "error"
			} else {
				checks["database"] = map[string]interface{}{
					"status": "ok",
				}
			}
		}
	} else {
		checks["database"] = map[string]interface{}{
			"status": "error",
			"error":  "Database not initialized",
		}
		overallStatus = "error"
	}
	
	// Configuration check
	cfg := config.Get()
	if cfg != nil {
		checks["config"] = map[string]interface{}{
			"status": "ok",
		}
	} else {
		checks["config"] = map[string]interface{}{
			"status": "error",
			"error":  "Configuration not loaded",
		}
		overallStatus = "error"
	}
	
	statusCode := fiber.StatusOK
	if overallStatus == "error" {
		statusCode = fiber.StatusServiceUnavailable
	}
	
	return c.Status(statusCode).JSON(fiber.Map{
		"status":    overallStatus,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
		"version":   "1.0.0",
		"service":   "go-api",
	})
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(startTime).String(),
		"version":   "1.0.0",
		"service":   "go-api",
	})
}

// MetricsCheck provides basic metrics
func (h *HealthHandler) MetricsCheck(c *fiber.Ctx) error {
	var dbStats sql.DBStats
	var dbConnections map[string]interface{}
	
	if h.db != nil {
		if sqlDB, err := h.db.DB(); err == nil {
			dbStats = sqlDB.Stats()
			dbConnections = map[string]interface{}{
				"open_connections":     dbStats.OpenConnections,
				"in_use":              dbStats.InUse,
				"idle":                dbStats.Idle,
				"wait_count":          dbStats.WaitCount,
				"wait_duration":       dbStats.WaitDuration.String(),
				"max_idle_closed":     dbStats.MaxIdleClosed,
				"max_lifetime_closed": dbStats.MaxLifetimeClosed,
			}
		} else {
			dbConnections = map[string]interface{}{
				"error": "Failed to get database stats",
			}
		}
	} else {
		dbConnections = map[string]interface{}{
			"error": "Database not initialized",
		}
	}
	
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"uptime":    time.Since(startTime).String(),
		"database":  dbConnections,
		"version":   "1.0.0",
		"service":   "go-api",
	})
}

// DetailedHealthCheck provides comprehensive health information
func (h *HealthHandler) DetailedHealthCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	
	checks := make(map[string]interface{})
	overallStatus := "ok"
	
	// Database connectivity check
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			checks["database"] = map[string]interface{}{
				"status": "error",
				"error":  "Failed to get database instance",
			}
			overallStatus = "error"
		} else {
			if err := sqlDB.PingContext(ctx); err != nil {
				checks["database"] = map[string]interface{}{
					"status": "error",
					"error":  err.Error(),
				}
				overallStatus = "error"
			} else {
				// Get database stats
				stats := sqlDB.Stats()
				checks["database"] = map[string]interface{}{
					"status": "ok",
					"connections": map[string]interface{}{
						"open":     stats.OpenConnections,
						"in_use":   stats.InUse,
						"idle":     stats.Idle,
					},
				}
			}
		}
	} else {
		checks["database"] = map[string]interface{}{
			"status": "error",
			"error":  "Database not initialized",
		}
		overallStatus = "error"
	}
	
	// Configuration check
	cfg := config.Get()
	if cfg != nil {
		checks["config"] = map[string]interface{}{
			"status": "ok",
			"environment": cfg.Environment,
		}
	} else {
		checks["config"] = map[string]interface{}{
			"status": "error",
			"error":  "Configuration not loaded",
		}
		overallStatus = "error"
	}
	
	statusCode := fiber.StatusOK
	if overallStatus == "error" {
		statusCode = fiber.StatusServiceUnavailable
	}
	
	return c.Status(statusCode).JSON(fiber.Map{
		"status":    overallStatus,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
		"uptime":    time.Since(startTime).String(),
		"version":   "1.0.0",
		"service":   "go-api",
	})
}

var startTime = time.Now()
