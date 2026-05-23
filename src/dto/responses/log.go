package responses

import (
	"api-log/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

// Log is the API response model for logs.
type Log struct {
	ID          uint      `json:"id"`
	Level       string    `json:"level"`
	Environment string    `json:"environment"`
	AppName     *string   `json:"appName"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	Route       *string   `json:"route"`
	Message     string    `json:"message"`
	Exception   *string   `json:"exception"`
	IpAddress   *string   `json:"ipAddress"`
}

// SetLog maps a database log model to an API response DTO.
func (l *Log) SetLog(log *models.Log) {
	l.ID = log.ID
	l.Level = log.Level
	l.Environment = log.Environment
	l.AppName = utils.PtrFromNullString(log.AppName)
	l.Version = log.Version
	l.CreatedAt = log.CreatedAt
	l.Route = utils.PtrFromNullString(log.Route)
	l.Message = log.Message
	l.Exception = utils.PtrFromNullString(log.Exception)
	l.IpAddress = utils.PtrFromNullString(log.IpAddress)
}
