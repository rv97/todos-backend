package kafka

import (
	"common/models"
)
type TodoCreatedMessage struct {
	models.TodoItem
}