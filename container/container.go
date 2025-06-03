package container

import (
	"context"
	"fmt"
	"go-api/app/repository"
	"go-api/app/service"
	"go-api/config"
	"go-api/database"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	Config       *config.Config
	DB           *gorm.DB
	AuthService  *service.AuthService
	EmailService *service.EmailService
}

func NewContainer() (*Container, error) {
	cfg := config.Get()
	db, err := database.InitDB()

	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	accessTokenRepo := repository.NewAccessTokenRepository(db)
	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, accessTokenRepo)
	emailService := service.NewEmailService(cfg)

	container := &Container{
		Config:       cfg,
		DB:           db,
		AuthService:  authService,
		EmailService: emailService,
	}

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
