package errors

import "github.com/gofiber/fiber/v2"

// Define error codes as constants.
const (
	NotFound     = "notFound"
	BodyParse    = "bodyParse"
	Validator    = "validator"
	CreateLog    = "createLog"
	Unauthorized = "unauthorized"
	// Add more error codes as needed.
)

// ErrorResponse creates a JSON response with a message and code.
func ErrorResponse(c *fiber.Ctx, status int, code, message interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"code":    code,
		"message": message,
	})
}
