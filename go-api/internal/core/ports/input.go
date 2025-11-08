package ports

import "southpark-api/internal/core/domain"

// MessageService defines the use cases for message handling
type MessageService interface {
	// SendMessage sends a message to the messaging system
	SendMessage(author, body string) error

	// CreateMessage creates a new message with the given author and body
	CreateMessage(author, body string) (*domain.Message, error)
}
