package database

import (
	"api-log/main/src/models"
	"context"
	"errors"
	"fmt"
	"time"

	utilsdb "github.com/ArnoldPMolenaar/api-utils/database"
	"gorm.io/gorm"
)

var Pg *gorm.DB

// OpenDBConnection Start a new database connection.
// Also tries to migrate the database schema.
func OpenDBConnection() error {
	// Open connection to database.
	db, err := utilsdb.PostgresSQLConnection()
	if err != nil {
		return err
	}

	// Migrate the database schema.
	err = Migrate(db)
	if err != nil {
		return err
	}

	// Set the global DB variable.
	Pg = db

	return nil
}

// ReadinessCheck verifies that the database connection is initialized and reachable.
func ReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	sqlDB, err := Pg.DB()
	if err != nil {
		return fmt.Errorf("database sql handle unavailable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// MigrationReadinessCheck verifies that required tables and seed data exist.
func MigrationReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	requiredTables := []any{&models.App{}, &models.LogLevel{}, &models.Log{}}
	for _, table := range requiredTables {
		if !Pg.Migrator().HasTable(table) {
			return fmt.Errorf("missing required table for %T", table)
		}
	}

	requiredLogLevels := []string{"Debug", "Info", "Success", "Warning", "Error", "Panic"}
	var count int64
	if err := Pg.Model(&models.LogLevel{}).Where("name IN ?", requiredLogLevels).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to validate log level seeds: %w", err)
	}

	if count < int64(len(requiredLogLevels)) {
		return fmt.Errorf("log level seeds incomplete: have %d, want %d", count, len(requiredLogLevels))
	}

	return nil
}
