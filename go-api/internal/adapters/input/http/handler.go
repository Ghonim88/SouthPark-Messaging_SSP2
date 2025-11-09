package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"southpark-api/internal/core/ports"
)

// MessageRequest represents the incoming JSON request
type MessageRequest struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

// MessageResponse represents the API response
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// MessageHandler handles HTTP requests for messages
type MessageHandler struct {
	service ports.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(service ports.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

// PostMessage handles POST /messages endpoint
func (h *MessageHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only accept POST method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(MessageResponse{
			Success: false,
			Error:   "Method not allowed, use POST",
		})
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var req MessageRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{
			Success: false,
			Error:   "Invalid JSON",
		})
		return
	}

	// Validate request
	if req.Author == "" || req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{
			Success: false,
			Error:   "Both 'author' and 'body' fields are required",
		})
		return
	}

	// Send message using the service (service handles creation/validation/publishing)
	if err := h.service.SendMessage(req.Author, req.Body); err != nil {
		log.Printf("Error sending message: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(MessageResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to send message: %v", err),
		})
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MessageResponse{
		Success: true,
		Message: "Message published successfully",
	})
	log.Printf("Message sent: Author=%s, Body=%s", req.Author, req.Body)
}

// HealthCheck handles GET /health endpoint
func (h *MessageHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "southpark-api",
	})
}
