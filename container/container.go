package container

import (
	"context"
	"fmt"
	"go-api/config"
	"go-api/database"
	authService "go-api/domain/auth/service"
	"go-api/email"
	"go-api/repository"
	"log"
	"time"

	"gorm.io/gorm"
)

type Container struct {
	Config       *config.Config
	DB           *gorm.DB
	AuthService  *authService.AuthService
	EmailService *email.EmailService
}

func NewContainer() (*Container, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("configuration not initialized")
	}
	
	// Initialize database with retry logic
	var db *gorm.DB
	var err error
	
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = database.InitDB()
		if err == nil {
			break
		}
		
		if i < maxRetries-1 {
			log.Printf("Database connection attempt %d failed: %v, retrying...", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database after %d attempts: %w", maxRetries, err)
	}

	// Initialize repositories with error handling
	userRepo := repository.NewUserRepository(db)
	if userRepo == nil {
		return nil, fmt.Errorf("failed to initialize user repository")
	}
	
	roleRepo := repository.NewRoleRepository(db)
	if roleRepo == nil {
		return nil, fmt.Errorf("failed to initialize role repository")
	}
	
	accessTokenRepo := repository.NewAccessTokenRepository(db)
	if accessTokenRepo == nil {
		return nil, fmt.Errorf("failed to initialize access token repository")
	}

	// Initialize services with error handling
	authService := authService.NewAuthService(userRepo, roleRepo, accessTokenRepo)
	if authService == nil {
		return nil, fmt.Errorf("failed to initialize auth service")
	}
	
	emailService := email.NewEmailClient(cfg)
	if emailService == nil {
		return nil, fmt.Errorf("failed to initialize email service")
	}

	container := &Container{
		Config:       cfg,
		DB:           db,
		AuthService:  authService,
		EmailService: emailService,
	}

	log.Println("Container initialized successfully")
	return container, nil
}

// Close gracefully shuts down the container and cleans up resources
func (c *Container) Close(ctx context.Context) error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}

		log.Println("Closing database connections...")
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}

	log.Println("Container cleanup completed")
	return nil
}
