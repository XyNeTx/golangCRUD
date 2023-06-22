package routes

import (
	"uplevel-api/configs"
	"uplevel-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	config := configs.AppConfig
	app.Get(config.APIURL+"/verify", controllers.GetValidToken)
	app.Get(config.APIURL+"/user/:userId", controllers.GetUserByID)
	app.Post(config.APIURL+"/user", controllers.CreateUser)
}
