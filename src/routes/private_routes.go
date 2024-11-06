package routes

import (
	"api-log/main/src/controllers"
	"api-log/main/src/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create protected routes group.
	route := a.Group("/v1", middleware.MachineProtected())

	// Routes for GET method:
	route.Get("/log", controllers.GetLogs)

	// Routes for POST method:
	route.Post("/log", controllers.CreateLog)
}
