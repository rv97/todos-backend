package main

import (
	"common/models"
	"log"
	"todos-api/repository"
	"todos-api/services"

	"github.com/gofiber/fiber/v2"
)

type TodoController struct {
	todoRepo      *repository.TodoRepo
	todoAnalytics *services.TodoAnalytics
}

type TodoCreateRequest struct {
	models.TodoItem
}

func NewTodoController(todoRepo *repository.TodoRepo, todoAnalytics *services.TodoAnalytics) *TodoController {
	return &TodoController{
		todoRepo,
		todoAnalytics,
	}
}

func (t *TodoController) addTodos(c *fiber.Ctx) error {
	todoCR := new(TodoCreateRequest)

	if err := c.BodyParser(todoCR); err != nil {
		log.Println(err)
	}

	t.todoRepo.CreateTodo(c.UserContext(), todoCR.TodoItem)
	go t.todoAnalytics.PublishCreatedTodo(&todoCR.TodoItem)
	c.Status(200).JSON(todoCR)
	return nil
}
