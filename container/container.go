package container

import (
	"context"
	"fmt"
	"go-api/config"
	"go-api/database"
	"go-api/internal/repository"
	"go-api/internal/service"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	Config      *config.Config
	DB          *gorm.DB
	AuthService *service.AuthService
	UserService *service.UserService
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

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo)
	userService := service.NewUserService(userRepo)

	container := &Container{
		Config:      cfg,
		DB:          db,
		AuthService: authService,
		UserService: userService,
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
