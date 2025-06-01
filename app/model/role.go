package model

type Role struct {
	BaseModelAttributes
	Code  string `gorm:"uniqueIndex;not null" json:"code"`
	Name  string `gorm:"not null" json:"name"`

	Users []User `gorm:"foreignKey:RoleID" json:"users"`
}
