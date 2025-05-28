package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Code  string `gorm:"uniqueIndex;not null" json:"code"`
	Name  string `gorm:"not null" json:"name"`

	Users []User `gorm:"foreignKey:RoleID" json:"users"`
}

