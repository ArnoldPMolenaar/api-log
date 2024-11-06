package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents the structure of the error response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler is a custom error handler for Fiber.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error.
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error.
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	} else {
		// If it's a panic error, include the panic message.
		message = err.Error()
	}

	// Create the error response.
	errorResponse := ErrorResponse{
		Code:    "internalServerError",
		Message: message,
	}

	// Return the error response as JSON.
	return c.Status(code).JSON(errorResponse)
}
