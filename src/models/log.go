package models

import (
	"time"
)

type Log struct {
	ID          uint   `gorm:"primaryKey"`
	Level       string `gorm:"not null"`
	Environment string `gorm:"not null"`
	Version     string `gorm:"not null"`
	CreatedAt   time.Time
	Route       string `gorm:"null"`
	Message     string `gorm:"not null"`
	Exception   string `gorm:"null"`
	IpAddress   string `gorm:"null"`

	LogLevel LogLevel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:Level;references:Name"`
}
