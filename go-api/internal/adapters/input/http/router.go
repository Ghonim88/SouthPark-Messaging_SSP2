package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux" // Using gorilla/mux for routing
)

// SetupRouter sets up the HTTP routes
func SetupRouter(handler *MessageHandler) *mux.Router {
	router := mux.NewRouter()

	// Define the POST /messages endpoint
	router.HandleFunc("/messages", handler.PostMessage).Methods("POST")
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Log all registered routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("Registered route: %s %s", methods, pathTemplate)
		}
		return nil
	})

	return router
}

// CORS Middleware to handle CORS (optional - for browser access)
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
