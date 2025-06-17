package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"go-api/model"
	"go-api/shared/timezone"

	"gorm.io/gorm"
)

type AccessTokenRepository struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) *AccessTokenRepository {
	return &AccessTokenRepository{db: db}
}

// generateSecureToken generates a cryptographically secure random token
func (r *AccessTokenRepository) generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (r *AccessTokenRepository) Create(ctx context.Context, userID uint, expiresIn time.Duration) (*model.AccessToken, error) {
	token, err := r.generateSecureToken()
	if err != nil {
		return nil, err
	}

	accessToken := &model.AccessToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: timezone.Now().Add(expiresIn), // Use timezone-aware time
	}

	if err := r.db.WithContext(ctx).Create(accessToken).Error; err != nil {
		return nil, err
	}

	// Load the user relationship
	if err := r.db.WithContext(ctx).Preload("User").First(accessToken, accessToken.ID).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

func (r *AccessTokenRepository) FindByToken(ctx context.Context, token string) (*model.AccessToken, error) {
	var accessToken model.AccessToken

	// Only find tokens that are not deleted
	err := r.db.WithContext(ctx).Preload("User").Where("token = ?", token).First(&accessToken).Error
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (r *AccessTokenRepository) RevokeToken(ctx context.Context, token string) error {
	// Soft delete the token by setting deleted_at
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&model.AccessToken{}).Error
}

func (r *AccessTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	// Soft delete all tokens for a user
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.AccessToken{}).Error
}

func (r *AccessTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	// Hard delete expired tokens that are already soft deleted
	return r.db.WithContext(ctx).Unscoped().Where("expires_at < ? AND deleted_at IS NOT NULL", timezone.Now()).Delete(&model.AccessToken{}).Error
}

// CleanupExpiredTokens deletes all expired tokens (even if not revoked)
// This should be called periodically to clean up the database
func (r *AccessTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).Unscoped().Where("expires_at < ?", timezone.Now()).Delete(&model.AccessToken{}).Error
}
