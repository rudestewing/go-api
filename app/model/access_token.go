package model

import (
	"go-api/app/shared/timezone"
	"time"
)

type AccessToken struct {
	BaseModelAttributes
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

// IsValid checks if the token is valid (not expired and not deleted)
func (at *AccessToken) IsValid() bool {
	return at.DeletedAt.Time.IsZero() && timezone.Now().Before(at.ExpiresAt)
}
