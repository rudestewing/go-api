package router

import (
	"go-api/app/handler"
	"go-api/container"
	"go-api/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *handler.Handler, container *container.Container) {
	router := app.Group("/api/v1")
	router.Use(middleware.TimeoutMiddleware(30*time.Second, "Operation timed out"))
	router.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data":    nil,
			"message": "API is running",
		})
	})

	auth := router.Group("/auth")
	auth.Post("/login", handler.AuthHandler.Login)
	auth.Post("/register", handler.AuthHandler.Register)

	protectedAuth :=  auth.Use(middleware.AuthMiddleware(container.AuthService))
	protectedAuth.Post("/logout", handler.AuthHandler.Logout)
	protectedAuth.Post("logout-all", handler.AuthHandler.LogoutAll)
}
