package main

import (
	"context"
	"github.com/gofiber/contrib/otelfiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"kafka"
	"log"
	"os"
	"todos-api/database"
	"todos-api/repository"
	"todos-api/services"

	"github.com/gofiber/fiber/v2"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// var Tracer = otel.Tracer("todos-api")
var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

func initTracer() *sdktrace.TracerProvider {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("todos-api"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	dsn := "host=localhost user=postgres password=password port=5432 sslmode=disable"
	db := database.NewDB(dsn)

	k := kafka.NewKafka("localhost:9092")

	todoAnalytics := services.NewTodoAnalytics(k)

	database.Migrate(db)

	todoRepo := repository.NewTodoRepo(db)
	todoController := NewTodoController(todoRepo, todoAnalytics)

	app := fiber.New()

	app.Use(otelfiber.Middleware())

	app.Post("/todos", todoController.addTodos)

	app.Listen(":3000")
}
