package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseAttributes struct {
	ID				  uint   `gorm:"primaryKey" json:"id"`
	CreatedAt		time.Time `json:"created_at"`
	UpdatedAt		time.Time `json:"updated_at"`
	DeletedAt		gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Use pointer to handle soft deletes
}
