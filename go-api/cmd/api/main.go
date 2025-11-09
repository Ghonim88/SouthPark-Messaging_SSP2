package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"southpark-api/config"
	httpAdapter "southpark-api/internal/adapters/input/http"
	"southpark-api/internal/adapters/output/rabbitmq"
	"southpark-api/internal/core/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	cfg.Print()

	// Initialize RabbitMQ publisher (Output Adapter)
	log.Println("Connecting to RabbitMQ...")
	publisher, err := rabbitmq.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

	// Initialize Message Service (Core Service)
	messageService := services.NewMessageService(publisher)

	// Initialize HTTP Handler (Input Adapter)
	messageHandler := httpAdapter.NewMessageHandler(messageService)

	// Setup router
	router := httpAdapter.SetupRouter(messageHandler)

	// Enable CORS middleware
	routerWithCORS := httpAdapter.EnableCORS(router)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: routerWithCORS,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s...", cfg.ServerPort)
		log.Println("Endpoints:")
		log.Printf("  POST http://localhost:%s/messages", cfg.ServerPort)
		log.Printf("  GET http://localhost:%s/health", cfg.ServerPort)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %s: %v", cfg.ServerPort, err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
