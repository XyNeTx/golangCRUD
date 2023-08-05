package main

import (
	"uplevel-api/configs"
	"uplevel-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config := configs.AppConfig
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	})

	//run database
	configs.ConnectDB()

	// routes
	routes.UserRoute(app)
	routes.OrganizationRoute(app)

	app.Listen(config.ServerAddress)
}
