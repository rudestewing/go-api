package handler

import (
	"go-api/shared/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// Tambahkan dependencies jika diperlukan
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id")
	email := c.Locals("email")

	userData := fiber.Map{
		"user_id": userId,
		"email":   email,
	}

	return response.Success(c, userData)
}
