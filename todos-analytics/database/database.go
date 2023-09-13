package database

import (
	"todos-analytics/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.TodoAnalyticsDto{})
}