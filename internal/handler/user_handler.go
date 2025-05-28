package handler

import (
	"go-api/internal/service"
	"go-api/shared/context"
	"go-api/shared/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile menampilkan profil user saat ini
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Create a context from the Fiber context
	ctx := context.CreateContext(c)
	
	// Extract user information from the Go context
	userId := context.GetUserID(ctx)
	email := context.GetEmail(ctx)

	userData := fiber.Map{
		"user_id": userId,
		"email":   email,
	}

	return response.Success(c, userData)
}

// GetUser mengambil data user berdasarkan ID dengan timeout default (30 detik)
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	// Create a context from the Fiber context
	baseCtx := context.CreateContext(c)
	
	// Gunakan GenerateContextWithTimeout dengan nilai default (0)
	ctx, cancel := context.GenerateContextWithTimeout(0, baseCtx)
	defer cancel() // Pastikan untuk memanggil cancel function ketika selesai
	
	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}
	
	return response.Success(c, user)
}

// GetAllUsers mengambil semua data user dengan timeout yang dapat dikonfigurasi
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Parse timeout parameter jika ada
	timeout := c.QueryInt("timeout", 0) // Default 0 akan menggunakan default timeout (30 detik)
	
	// Create a context from the Fiber context
	baseCtx := context.CreateContext(c)
	
	// Gunakan GenerateContextWithTimeout dengan nilai timeout yang ditentukan
	// Jika timeout > 0, gunakan nilai tersebut, jika tidak akan menggunakan default 30 detik
	ctx, cancel := context.GenerateContextWithTimeout(time.Duration(timeout)*time.Second, baseCtx)
	defer cancel() // Pastikan untuk memanggil cancel function ketika selesai
	
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve users")
	}
	
	return response.Success(c, users)
}
