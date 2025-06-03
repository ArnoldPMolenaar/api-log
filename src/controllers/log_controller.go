package controllers

import (
	"api-log/main/src/database"
	"api-log/main/src/enums"
	"api-log/main/src/errors"
	"api-log/main/src/models"
	"api-log/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateLog function creates a new log in the database.
func CreateLog(c *fiber.Ctx) error {
	log := &models.Log{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(log); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.BodyParse, err.Error())
	}

	if (log.Level == enums.Error.String() || log.Level == enums.Panic.String()) && *log.Exception == "" {
		return errorutil.Response(
			c,
			fiber.StatusBadRequest,
			errors.Validator,
			"Exception field is required for error and panic logs.",
		)
	}

	// Validate log fields.
	validate := utils.NewValidator()
	if err := validate.Struct(log); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.Validator, utils.ValidatorErrors(err))
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(*log.AppName); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppExists, "AppName does not exist.")
	}

	// Save log to database.
	result := database.Pg.Create(&log)
	if result.Error != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.CreateLog, result.Error.Error())
	}

	// Return HTTP 200 status and JSON response.
	return c.Status(fiber.StatusOK).JSON(log)
}

// GetLogs function fetches all logs from the database.
func GetLogs(c *fiber.Ctx) error {
	logs := make([]models.Log, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":          true,
		"level":       true,
		"environment": true,
		"app_name":    true,
		"version":     true,
		"created_at":  true,
		"route":       true,
		"message":     true,
		"exception":   true,
		"ip_address":  true,
	}

	queryFunc := pagination.Query(values, allowedColumns)
	sortFunc := pagination.Sort(values, allowedColumns)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)

	db := database.Pg.Scopes(queryFunc, sortFunc).Limit(limit).Offset(offset).Find(&logs)
	if db.Error != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.GetLogs, db.Error.Error())
	}

	total := int64(0)
	database.Pg.Scopes(queryFunc).Model(&models.Log{}).Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), logs)

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}
