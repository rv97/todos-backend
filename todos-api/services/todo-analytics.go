package services

import (
	"common/models"
	"encoding/json"
	"kafka"
	"log"
)

type TodoAnalytics struct {
	kafka *kafka.Kafka
}

func NewTodoAnalytics(kafka *kafka.Kafka) *TodoAnalytics {
	return &TodoAnalytics{
		kafka,
	}
}

func (t *TodoAnalytics) PublishCreatedTodo(todoItem *models.TodoItem) {
	todoCreatedMessage := kafka.TodoCreatedMessage {
		TodoItem: *todoItem,
	}

	jsonValue, err := json.Marshal(todoCreatedMessage)

	if err != nil {
		log.Panic("Err: Cannot serialize JSON %w", err)
	}
	t.kafka.Produce(kafka.TodoCreated, todoItem.Title, jsonValue)
}