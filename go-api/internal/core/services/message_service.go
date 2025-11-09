package services

import (
	"fmt"
	"log"
	"southpark-api/internal/core/domain"
	"southpark-api/internal/core/ports"
)

// messageService implements the MessageService interface
type messageService struct {
	publisher ports.MessagePublisher
}

// NewMessageService creates a new instance of messageService
// This follows Dependency Injection - we inject the MessagePublisher (adapter)
// into the service, allowing for loose coupling and easier testing. and keep service independent from rabbitmq.
func NewMessageService(publisher ports.MessagePublisher) ports.MessageService {
	return &messageService{
		publisher: publisher,
	}
}

// CreateMessage creates and validates a new message
func (s *messageService) CreateMessage(author, body string) (*domain.Message, error) {
	message, err := domain.NewMessage(author, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	log.Printf("Message created: Author=%s, Body=%s", author, body)
	return message, nil
}

// SendMessage sends a message through the publisher
func (s *messageService) SendMessage(author, body string) error {
	//Create and validate the message
	message, err := s.CreateMessage(author, body)
	if err != nil {
		return err
	}

	//Publish the message using the injected publisher
	if err := s.publisher.Publish(message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	log.Printf("Message sent: Author=%s, Body=%s", author, body)
	return nil
}
