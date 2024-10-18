package models

import "gorm.io/gorm"

type LogLevel struct {
	gorm.Model
	Name string `gorm:"primaryKey"`
}
