package controllers

import "github.com/gofiber/fiber/v2"

func CreateLog(c *fiber.Ctx) error {
	// Return HTTP 200 status and JSON response.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "log created",
	})
}
