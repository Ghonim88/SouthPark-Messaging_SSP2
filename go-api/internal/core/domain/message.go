package domain

import (
	"encoding/json"
	"errors"
	"time"
)

// Message represents a South Park character message
type Message struct {
	Author string    `json:"author"`
	Body   string    `json:"body"`
	SentAt time.Time `json:"sent_at"`
}

// NewMessage creates a new Message instance with validation
func NewMessage(author, body string) (*Message, error) {
	if author == "" {
		return nil, errors.New("author is required")
	}
	if body == "" {
		return nil, errors.New("body is required")
	}
	return &Message{
		Author: author,
		Body:   body,
		SentAt: time.Now(),
	}, nil
}

// ToJSON converts message to Json bytes
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// Validate checks if the message has valid fields
func (m *Message) Validate() error {
	if m.Author == "" {
		return errors.New("author is required")
	}
	if m.Body == "" {
		return errors.New("body is required")
	}
	return nil
}
