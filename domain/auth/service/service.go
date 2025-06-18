package service

import (
	"context"
	"errors"
	"go-api/app"
	"go-api/config"
	entity "go-api/domain/auth/entity"
	"go-api/model"
	"go-api/repository"
	"go-api/shared/constant"
	"go-api/shared/logger"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	provider 			*app.Provider
	userRepo        *repository.UserRepository
	roleRepo        *repository.RoleRepository
	accessTokenRepo *repository.AccessTokenRepository
	tokenExpiry     time.Duration
}

func NewAuthService(p *app.Provider) *AuthService {
	return &AuthService{
		provider: 			 p,
		userRepo:        repository.NewUserRepository(p.DB),
		roleRepo:        repository.NewRoleRepository(p.DB),
		accessTokenRepo: repository.NewAccessTokenRepository(p.DB),
		tokenExpiry:     config.Get().JWTExpiry,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.AccessToken, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Create access token
	accessToken, err := s.accessTokenRepo.Create(ctx, user.ID, s.tokenExpiry)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*model.AccessToken, error) {
	accessToken, err := s.accessTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if !accessToken.IsValid() {
		return nil, errors.New("invalid or expired token")
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.accessTokenRepo.RevokeToken(ctx, token)
}

func (s *AuthService) LogoutAll(ctx context.Context, userID uint) error {
	return s.accessTokenRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *AuthService) Register(ctx context.Context, req *entity.RegisterRequest) error {
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

	go func() {
		if err := s.provider.Email.SendWelcomeEmail(user.Email, user.Name); err != nil {
			logger.Errorf("Failed to send email: %v", err)
		}
	}()

	return s.userRepo.Create(ctx, user)
}
