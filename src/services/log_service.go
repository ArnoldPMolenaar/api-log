package services

import (
	"api-log/main/src/database"
	"api-log/main/src/dto/requests"
	"api-log/main/src/dto/responses"
	"api-log/main/src/models"
	"strconv"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/valyala/fasthttp"
)

var allowedLogColumns = map[string]bool{
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

// GetLogs returns paginated logs as response DTOs.
func GetLogs(values *fasthttp.Args, pageValue, limitValue string) (pagination.Model, error) {
	logs := make([]models.Log, 0)

	queryFunc := pagination.Query(values, allowedLogColumns)
	sortFunc := pagination.Sort(values, allowedLogColumns)
	page, err := strconv.Atoi(pageValue)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitValue)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)

	db := database.Pg.Scopes(queryFunc, sortFunc).Limit(limit).Offset(offset).Find(&logs)
	if db.Error != nil {
		return pagination.Model{}, db.Error
	}

	total := int64(0)
	countQuery := database.Pg.Scopes(queryFunc).Model(&models.Log{}).Count(&total)
	if countQuery.Error != nil {
		return pagination.Model{}, countQuery.Error
	}
	pageCount := pagination.Count(int(total), limit)

	responseLogs := make([]responses.Log, 0, len(logs))
	for i := range logs {
		responseLog := responses.Log{}
		responseLog.SetLog(&logs[i])
		responseLogs = append(responseLogs, responseLog)
	}

	return pagination.CreatePaginationModel(limit, page, pageCount, int(total), responseLogs), nil
}

// CreateLog saves a new log in the database.
func CreateLog(request *requests.CreateLog) (*models.Log, error) {
	log := request.ToModel()

	if err := database.Pg.Create(log).Error; err != nil {
		return nil, err
	}

	return log, nil
}
