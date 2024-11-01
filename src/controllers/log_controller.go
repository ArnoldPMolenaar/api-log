package controllers

import (
	"api-log/main/src/database"
	"api-log/main/src/enums"
	"api-log/main/src/errors"
	"api-log/main/src/models"
	"api-log/main/src/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateLog(c *fiber.Ctx) error {
	log := &models.Log{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(log); err != nil {
		return errors.ErrorResponse(c, fiber.StatusBadRequest, errors.BodyParse, err.Error())
	}

	if (log.Level == enums.Error.String() || log.Level == enums.Panic.String()) && log.Exception == "" {
		return errors.ErrorResponse(c, fiber.StatusBadRequest, errors.Validator, "Exception field is required for error and panic logs.")
	}

	// Validate log fields.
	validate := utils.NewValidator()
	if err := validate.Struct(log); err != nil {
		return errors.ErrorResponse(c, fiber.StatusBadRequest, errors.Validator, utils.ValidatorErrors(err))
	}

	// Save log to database.
	result := database.Pg.Create(&log)
	if result.Error != nil {
		return errors.ErrorResponse(c, fiber.StatusInternalServerError, errors.CreateLog, result.Error.Error())
	}

	// Return HTTP 200 status and JSON response.
	return c.Status(fiber.StatusOK).JSON(log)
}
