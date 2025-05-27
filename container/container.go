package container

import (
	"fmt"
	"go-api/internal/handler"
	"go-api/internal/repository"
	"go-api/internal/service"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
	AppPort     string
	JWTSecret   string
}

type Container struct {
	Config      *Config
	DB          *gorm.DB
	UserRepo    *repository.UserRepository
	AuthService *service.AuthService
	AuthHandler *handler.AuthHandler
}

func loadConfig() *Config {
	if err := godotenv.Load(); err != nil {
	log.Println("Warning: Error loading .env file:", err)
}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AppPort:    os.Getenv("APP_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}


func initDB(cfg *Config) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the user model
	// if err := db.AutoMigrate(&model.User{}); err != nil {
	// 	return nil, fmt.Errorf("failed to migrate database: %w", err)
	// }

	log.Println("Database connection established")
	return db, nil
}


func NewContainer() (*Container, error) {
	cfg := loadConfig()

	db, err := initDB(cfg)
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
