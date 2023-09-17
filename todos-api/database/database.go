package database

import (
	"gorm.io/plugin/opentelemetry/tracing"
	"log"
	"todos-api/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panic("Err: Db not created %w", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&repository.TodoCreateDto{})

	if err != nil {
		panic(err)
	}
}
