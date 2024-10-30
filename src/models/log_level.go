package models

type LogLevel struct {
	Name string `gorm:"uniqueIndex:idx_name;primaryKey;not null;autoIncrement:false"`
}
