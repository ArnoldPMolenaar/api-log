package models

import (
	"database/sql"
	"time"
)

type Log struct {
	ID          uint           `gorm:"primaryKey"`
	Level       string         `gorm:"not null"`
	Environment string         `gorm:"not null"`
	AppName     sql.NullString `gorm:"null"`
	Version     string         `gorm:"not null"`
	CreatedAt   time.Time
	Route       sql.NullString `gorm:"null"`
	Message     string         `gorm:"not null"`
	Exception   sql.NullString `gorm:"null"`
	IpAddress   sql.NullString `gorm:"null"`

	LogLevel LogLevel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:Level;references:Name" json:"-"`
	App      *App     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name" json:"-"`
}
