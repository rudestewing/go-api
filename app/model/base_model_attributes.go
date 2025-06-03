package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModelAttributes struct {
	ID				  uint   `gorm:"primaryKey" json:"id"`
	CreatedAt		time.Time `json:"created_at"`
	UpdatedAt		time.Time `json:"updated_at"`
	DeletedAt		gorm.DeletedAt `gorm:"index" json:"deleted_at"` // Use pointer to handle soft deletes
}


// BeforeCreate hook - automatically called by GORM
func (b *BaseModelAttributes) BeforeCreate(tx *gorm.DB) error {
    b.CreatedAt = time.Now()
    b.UpdatedAt = time.Now()
    return nil
}

// BeforeUpdate hook - automatically called by GORM
func (b *BaseModelAttributes) BeforeUpdate(tx *gorm.DB) error {
    b.UpdatedAt = time.Now()
    return nil
}