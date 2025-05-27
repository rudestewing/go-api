package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Name     string `gorm:"nullable" json:"name"`
	Password string `gorm:"nullable" json:"-"`
}
