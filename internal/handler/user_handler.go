package handler

import (
	"go-api/internal/service"
	"go-api/shared/response"

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
	ctx := c.UserContext()

	userData := fiber.Map{
		"user": ctx.Value("user_id"),
	}

	return response.Success(c, userData)
}

// GetUser mengambil data user berdasarkan ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, user)
}

// GetAllUsers mengambil semua data user
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Gunakan context dari Fiber yang sudah memiliki timeout dari middleware
	ctx := c.UserContext()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve users")
	}

	return response.Success(c, users)
}
