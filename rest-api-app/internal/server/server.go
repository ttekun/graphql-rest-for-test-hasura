package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ttekun/rest-api-app/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	port string
}

// NewServer creates a new server instance
func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Customer endpoints
	http.HandleFunc("/Customer/getCustomerInfo", handlers.GetCustomerInfo)

	// Health check endpoint
	http.HandleFunc("/healthz", handlers.HealthCheck)

	// Start server
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, nil)
}