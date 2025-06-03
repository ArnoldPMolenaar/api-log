package models

import (
	"time"
)

type Log struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Level       string    `gorm:"not null" json:"level"`
	Environment string    `gorm:"not null" json:"environment"`
	AppName     string    `gorm:"null" json:"appName"`
	Version     string    `gorm:"not null" json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	Route       string    `gorm:"null" json:"route"`
	Message     string    `gorm:"not null" json:"message"`
	Exception   string    `gorm:"null" json:"exception"`
	IpAddress   string    `gorm:"null" json:"ipAddress"`

	LogLevel LogLevel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:Level;references:Name" json:"-"`
	App      App      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name" json:"-"`
}
