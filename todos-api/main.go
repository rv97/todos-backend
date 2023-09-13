package main

import (
	"kafka"
	"log"
	"todos-api/database"
	"todos-api/repository"
	"todos-api/services"

	"github.com/gofiber/fiber/v2"
)


func main() {
	dsn := "host=localhost user=postgres password=password port=5432 sslmode=disable"
	db, err := database.NewDB(dsn)
	if err != nil {
		log.Panic("Err: DB not created %w", err)
	}
	

	k := kafka.NewKafka("localhost:9092")

	todoAnalytics := services.NewTodoAnalytics(k)

	database.Migrate(db)

	todoRepo := repository.NewTodoRepo(db)
	todoController := NewTodoController(todoRepo, todoAnalytics)

	app := fiber.New()

	app.Post("/todos", todoController.addTodos)

	app.Listen(":3000")
}