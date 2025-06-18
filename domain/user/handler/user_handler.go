package handler

import (
	"strconv"

	"go-api/domain/user/service"
	"go-api/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserHandler demonstrates how to create handlers with proper dependency injection
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler with proper dependency injection
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(db),
	}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": user,
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.userService.CreateUser(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": user,
	})
}

// GetUserByEmail retrieves a user by email
func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email parameter is required",
		})
	}

	user, err := h.userService.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(fiber.Map{
		"data": user,
	})
}

// RegisterUserRoutes registers user routes using dependency injection
func RegisterUserRoutes(app *fiber.App, db *gorm.DB) {
	// Using instance-based handler with dependency injection
	handler := NewUserHandler(db)
	
	userGroup := app.Group("/api/users")
	userGroup.Get("/:id", handler.GetUser)
	userGroup.Post("/", handler.CreateUser)
	userGroup.Get("/", handler.GetUserByEmail)
}
