package dto

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

// RegisterRequest represents the user registration request payload
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100,alpha" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,strong_password" example:"Password123!"`
}

