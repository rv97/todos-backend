package repository

import (
	commonModels "common/models"
	"log"
	"todos-analytics/models"

	"gorm.io/gorm"
)

type TodoAnalytics struct {
	db *gorm.DB
}

func NewTodoAnalytics(db *gorm.DB) *TodoAnalytics {
	return &TodoAnalytics{db}
}

func (t *TodoAnalytics) CreateAnalyticsItemInDB(item commonModels.TodoItem) {
	var todoAnalyticsDto models.TodoAnalyticsDto
	if err := t.db.First(&todoAnalyticsDto).Error; err != nil {
		log.Panic("Err: Cannot read from table %w", err)
	}

	log.Printf("in CreateAnalytics: %v", todoAnalyticsDto.Count)
	todoAnalyticsDto.Count++

	if err := t.db.Save(&todoAnalyticsDto).Error; err != nil {
        log.Fatalf("Failed to update the record: %v", err)
    }
}


