package requests

import (
	"api-log/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

// CreateLog is the API request model for creating logs.
type CreateLog struct {
	Level       string  `json:"level"`
	Environment string  `json:"environment"`
	AppName     *string `json:"appName"`
	Version     string  `json:"version"`
	Route       *string `json:"route"`
	Message     string  `json:"message"`
	Exception   *string `json:"exception"`
	IpAddress   *string `json:"ipAddress"`
}

// ToModel maps the request DTO to the database model.
func (r *CreateLog) ToModel() *models.Log {
	return &models.Log{
		Level:       r.Level,
		Environment: r.Environment,
		AppName:     utils.NewNullString(r.AppName),
		Version:     r.Version,
		Route:       utils.NewNullString(r.Route),
		Message:     r.Message,
		Exception:   utils.NewNullString(r.Exception),
		IpAddress:   utils.NewNullString(r.IpAddress),
	}
}
