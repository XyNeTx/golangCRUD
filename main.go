package main

import (
	"log"
	"uplevel-api/configs"
	"uplevel-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB"})
	})

	//run database
	configs.ConnectDB()

	// routes
	routes.UserRoute(app)

	app.Listen(config.ServerAddress)
}
