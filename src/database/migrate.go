package database

import (
	"api-log/main/src/models"
	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.App{}, &models.LogLevel{}, &models.Log{})
	if err != nil {
		return err
	}

	// Seed log levels
	logLevels := []string{"Debug", "Info", "Success", "Warning", "Error", "Panic"}
	for _, level := range logLevels {
		var logLevel models.LogLevel
		if err := db.FirstOrCreate(&logLevel, models.LogLevel{Name: level}).Error; err != nil {
			return err
		}
	}

	return nil
}
