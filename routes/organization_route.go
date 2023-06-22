package routes

import (
	"uplevel-api/configs"
	"uplevel-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func OrganizationRoute(app *fiber.App) {
	config := configs.AppConfig
	app.Get(config.APIURL+"/:organizationId/summary", controllers.GetOrgSummary)
	app.Get(config.APIURL+"/organization/:userId", controllers.GetMyOrg)
	app.Get(config.APIURL+"/organization/", controllers.GetAllOrgs)
	app.Post(config.APIURL+"/organization/:userId", controllers.CreateOrg)
	app.Patch(config.APIURL+"/organization/:organizationId", controllers.EditOrg)
	app.Delete(config.APIURL+"/organization/:organizationId", controllers.DeleteOrg)
}
