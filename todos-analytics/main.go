package main

import (
	"common/models"
	"encoding/json"
	"kafka"
	"log"
	"sync"
	"todos-analytics/database"
	localModel "todos-analytics/models"
	"todos-analytics/repository"
)

func handlerClosureScope(todoAnalytics *repository.TodoAnalytics) func(kafka.Message) {
	return func(e kafka.Message) {
		todoCreatedMessage := new(kafka.TodoCreatedMessage)
		
		err := json.Unmarshal(e.Value, &todoCreatedMessage)
	
		if err != nil {
			log.Panic("Err: Cannot deserialize contents %w", err)
		}

		entryItem := &models.TodoItem{
			Title: todoCreatedMessage.Title,
			Description: todoCreatedMessage.Description,
			IsCompleted: todoCreatedMessage.IsCompleted,
		}
		log.Printf("Inside handler!")
		todoAnalytics.CreateAnalyticsItemInDB(*entryItem)
	}
}

func main() {
	dsn := "host=localhost user=postgres password=password port=5432 sslmode=disable"
	db, err := database.NewDB(dsn)

	if err != nil {
		log.Panic("Err: Database not created %w", err)
	}

	todoAnalytics := repository.NewTodoAnalytics(db)
	
	database.Migrate(db)

	db.Create(&localModel.TodoAnalyticsDto{Count: 0})

	k := kafka.NewKafka("localhost:9092")

	numConsumers := 3

	// Create a wait group to wait for all consumers to finish
	var wg sync.WaitGroup

	for i := 0; i < numConsumers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			k.Consume("todos-analytics", []kafka.Topics{kafka.TodoCreated}, handlerClosureScope(todoAnalytics))
		}()
	}

	wg.Wait()
}