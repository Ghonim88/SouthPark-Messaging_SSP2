package rabbitmq

import (
	"fmt"
	"log"
	"southpark-api/internal/core/domain"
	"southpark-api/internal/core/ports"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QUEUE_NAME = "southpark_messages"
)

// rabbitMQPublisher implements the MessagePublisher interface
type rabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher instance
func NewRabbitMQPublisher(rabbitURL string) (ports.MessagePublisher, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare the queue (idempotent -won't recreate if it exists)
	queue, err := channel.QueueDeclare(
		QUEUE_NAME,
		true,  // durable - survives broker restart
		false, // auto-delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}
	log.Printf("Connected to RabbitMQ, queue declared: %s", QUEUE_NAME)

	// Return the RabbitMQ publisher
	return &rabbitMQPublisher{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

// Publish publishes a message to RabbitMQ
func (p *rabbitMQPublisher) Publish(message *domain.Message) error {
	//Convert message to JSON
	body, err := message.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	//Publish to the queue
	err = p.channel.Publish(
		"",           // exchange (empty string = default exchange)
		p.queue.Name, // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // make message persistent
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	log.Printf("Message published to RabbitMQ: Author=%s, Body=%s", message.Author, message.Body)
	return nil

}

// Close closes the RabbitMQ connection and channel
func (p *rabbitMQPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	log.Println("RabbitMQ connection and channel closed")
	return nil
}
