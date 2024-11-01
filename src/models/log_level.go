package models

type LogLevel struct {
	Name string `gorm:"primaryKey:true;not null;autoIncrement:false"`
}
