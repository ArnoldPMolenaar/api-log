package controllers

import (
	"api-log/main/src/dto/requests"
	"api-log/main/src/dto/responses"
	"api-log/main/src/enums"
	"api-log/main/src/errors"
	"api-log/main/src/services"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"
)

// CreateLog function creates a new log in the database.
func CreateLog(c fiber.Ctx) error {
	request := requests.CreateLog{}

	// Check, if received JSON data is parsed.
	if err := c.Bind().Body(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err)
	}

	if (request.Level == enums.Error.String() || request.Level == enums.Panic.String()) && (request.Exception == nil || *request.Exception == "") {
		return errorutil.Response(
			c,
			fiber.StatusBadRequest,
			errorutil.Validator,
			"Exception field is required for error and panic logs.",
		)
	}

	// Validate log fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}

	// Check if app exists.
	if request.AppName != nil {
		if available, err := services.IsAppAvailable(*request.AppName); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
		}
	}

	log, err := services.CreateLog(&request)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.CreateLog, err.Error())
	}

	response := responses.Log{}
	response.SetLog(log)

	// Return HTTP 200 status and JSON response.
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetLogs function fetches all logs from the database.
func GetLogs(c fiber.Ctx) error {
	paginationModel, err := services.GetLogs(
		c.Request().URI().QueryArgs(),
		c.Query("page", "1"),
		c.Query("limit", "10"),
	)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.GetLogs, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}
