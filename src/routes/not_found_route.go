package routes

import (
	"api-log/main/src/errors"
	"github.com/gofiber/fiber/v2"
)

// NotFoundRoute func for describe 404 Error route.
func NotFoundRoute(a *fiber.App) {
	// Register new special route.
	a.Use(
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return errors.ErrorResponse(
				c,
				fiber.StatusNotFound,
				errors.NotFound,
				"sorry, endpoint is not found",
			)
		},
	)
}
