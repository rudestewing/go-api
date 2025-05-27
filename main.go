package main

import (
	"go-api/container"
	"go-api/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	container, err := container.NewContainer()

	if(err != nil) {
		panic(err.Error())
	}

	router.RegisterRoutes(app, container)

	app.Listen(":8000")
}
