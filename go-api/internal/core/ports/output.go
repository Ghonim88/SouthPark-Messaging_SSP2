package ports

import "southpark-api/internal/core/domain"

// MessagePublisher defines the interface for publishing messages
// This is the OUTPUT port - adapters will implement this
type MessagePublisher interface {
	// Publish publishes a message to an external system
	Publish(message *domain.Message) error

	// Close closes the connection to the message broker
	Close() error
}
