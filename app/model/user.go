package model

type User struct {
	BaseAttributes
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Name     string `gorm:"nullable" json:"name"`
	Password string `gorm:"nullable" json:"-"`
	RoleID   uint   `gorm:"not null" json:"role_id"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}
