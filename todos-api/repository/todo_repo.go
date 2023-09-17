package repository

import (
	"common/models"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

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

func (tr *TodoRepo) CreateTodo(ctx context.Context, todo models.TodoItem) {
	tracer := otel.GetTracerProvider().Tracer("db-repository")
	ctx, span := tracer.Start(ctx, "db-createTodo", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	span.AddEvent("creating a new model")
	newEntry := &TodoCreateDto{
		TodoItem: todo,
	}

	span.AddEvent("calling DB")
	tr.db.WithContext(ctx).Create(newEntry)
	span.SetAttributes(attribute.String("todo_title", todo.Title))
	//To track the status manually, do the following...
	//span.SetStatus(codes.Error, "Test Error")
}
