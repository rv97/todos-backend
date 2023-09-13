package repository

import (
	"common/models"

	"gorm.io/gorm"
)

type TodoRepo struct {
	db *gorm.DB
}

func NewTodoRepo(db *gorm.DB) *TodoRepo {
	return &TodoRepo{db}
}

type TodoCreateDto struct {
	gorm.Model
	models.TodoItem
}

func (tr *TodoRepo) CreateTodo (todo models.TodoItem) {
	newEntry := &TodoCreateDto{
		TodoItem: todo,
	}

	tr.db.Create(newEntry)
}