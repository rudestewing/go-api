package service

import (
	"context"
	"fmt"
	"go-api/model"
	"go-api/repository"

	"gorm.io/gorm"
)

// UserService demonstrates how to create a service with proper dependency injection
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service with proper dependency injection
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(db),
	}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}
	
	return s.userRepo.Create(ctx, user)
}


