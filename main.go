package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

  err := godotenv.Load()
	
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	// Example route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	// Start server
	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}