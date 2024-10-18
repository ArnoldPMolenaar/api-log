package models

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Level       string `gorm:"not null"`
	Environment string `gorm:"not null"`
	Version     string `gorm:"not null"`
	Route       string `gorm:"null"`
	Message     string `gorm:"not null"`
	Exception   string `gorm:"null"`

	LogLevel LogLevel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
