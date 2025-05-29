package handler

import (
	"go-api/app/dto"
	"go-api/app/service"
	"go-api/app/shared/response"
	"go-api/app/shared/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request
	if err := validator.ValidateStruct(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()

	// Call service with timeout context
	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}


	return response.Success(c, fiber.Map{
		"token": token,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate request structure
	if err := validator.ValidateStruct(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	// Additional password validation
	if err := validator.ValidatePassword(req.Password); err != nil {
		return response.ValidationError(c, err.Error())
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()

	// Call service with context
	if err := h.authService.Register(ctx, &req); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, nil, "User registered successfully")
}
