package routes

import (
	"log"
	"uplevel-api/configs"
	"uplevel-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func OrganizationRoute(app *fiber.App) {
	config, err := configs.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// app.Get(config.APIURL+"/:organizationId/summary", controllers.g)
	app.Get(config.APIURL+"/organization/:userId", controllers.GetMyOrg)
	app.Get(config.APIURL+"/organization/", controllers.GetAllOrgs)
	app.Post(config.APIURL+"/organization/", controllers.CreateOrg)
}
