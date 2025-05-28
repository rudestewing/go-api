package service

import (
	"context"
	"go-api/internal/model"
	"go-api/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}


func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}