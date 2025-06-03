package routes

import (
	"api-log/main/src/controllers"
	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create protected routes group.
	route := a.Group("/v1", middleware.MachineProtected())

	// Register route for /v1/apps.
	route.Post("/apps", controllers.CreateApp)

	// Routes for GET method:
	route.Get("/logs", controllers.GetLogs)

	// Routes for POST method:
	route.Post("/logs", controllers.CreateLog)
}
