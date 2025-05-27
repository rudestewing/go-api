package container

import (
	"context"
	"fmt"
	"go-api/config"
	"go-api/internal/handler"
	"go-api/internal/repository"
	"go-api/internal/service"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Container struct {
	Config      *config.Config
	DB          *gorm.DB
	UserRepo    *repository.UserRepository
	AuthService *service.AuthService
	AuthHandler *handler.AuthHandler
}

func initDB() (*gorm.DB, error) {
	cfg := config.Get()

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Configure connection pool to prevent memory leaks
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings from config
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)

	// Auto migrate the user model
	// if err := db.AutoMigrate(&model.User{}); err != nil {
	// 	return nil, fmt.Errorf("failed to migrate database: %w", err)
	// }

	log.Println("Database connection established")
	return db, nil
}


func NewContainer() (*Container, error) {
	cfg := config.Get()
	db, err := initDB()

	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	container := &Container{
		Config:      cfg,
		DB:          db,
		UserRepo:    userRepo,
		AuthService: authService,
		AuthHandler: authHandler,
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
