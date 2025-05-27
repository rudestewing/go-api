package dto

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents the user registration request payload
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

// User represents user data in responses
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Error  string              `json:"error"`
	Fields map[string][]string `json:"fields,omitempty"`
}
