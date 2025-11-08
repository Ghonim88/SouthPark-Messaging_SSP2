package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	ServerPort  string
	RabbitMQURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Print logs the Configuration
func (c *Config) Print() {
	fmt.Println("================================")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  ServerPort: %s\n", c.ServerPort)
	fmt.Printf("  RabbitMQURL: %s\n", c.RabbitMQURL)
	fmt.Println("================================")
}
