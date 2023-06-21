package routes

import (
	"log"
	"uplevel-api/configs"
	"uplevel-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	config, err := configs.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	app.Get(config.APIURL+"/user/:userId", controllers.GetUser)
	app.Post(config.APIURL+"/user", controllers.CreateUser)
}