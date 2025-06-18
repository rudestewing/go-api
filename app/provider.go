package app

import (
	"fmt"
	"go-api/config"
	"go-api/database"
	"go-api/email"
	"go-api/shared/logger"

	"gorm.io/gorm"
)

type Provider struct {
	DB    *gorm.DB
	Email *email.EmailService
}

func BootProvider(cfg *config.Config) (*Provider, error) {
	// Initialize database connection directly (not global)
	logger.Infof("Initializing database connection...")
	db, err := database.NewConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	logger.Infof("Database initialized successfully")
	logger.Infof("Initializing email service...")
	// Initialize email service as dependency
	emailService := email.NewEmailService(cfg)
	logger.Infof("Email service initialized successfully")

	return &Provider{
		DB:    db,
		Email: emailService,
	}, nil
}

func (p *Provider) ShutdownProvider() {
	if p.DB != nil {
		sqlDB, err := p.DB.DB()
		if err != nil {
			logger.Errorf("Error getting underlying sql.DB: %v", err)
			return
		}
		
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("Error closing database: %v", err)
		} else {
			logger.Infof("Database closed successfully")
		}
	}
}
