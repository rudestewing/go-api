package service

import (
	"context"
	"errors"
	"go-api/app/dto"
	"go-api/app/model"
	"go-api/app/repository"
	"go-api/app/shared/constant"
	"go-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	roleRepo    *repository.RoleRepository
	tokenExpiry time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		tokenExpiry: config.Get().JWTExpiry,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	})

	jwtSecret := config.Get().JWTSecret

	return token.SignedString([]byte(jwtSecret))
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	roleUser, err := s.roleRepo.FindByCode(ctx, constant.RoleCodeUser)

	if err != nil {
		return err
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		RoleID:   roleUser.ID,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(ctx, user)
}
