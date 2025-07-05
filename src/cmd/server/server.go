package main

import (
	"net/http"

	"github.com/wizenheimer/bloombox/pkg/emailchecker"
)

// Server represents the HTTP server for the email validation API
type Server struct {
	checker *emailchecker.EmailChecker
	config  *emailchecker.Config
}

// NewServer creates a new server instance with the given checker and config
func NewServer(checker *emailchecker.EmailChecker, config *emailchecker.Config) *Server {
	return &Server{
		checker: checker,
		config:  config,
	}
}

// SetupRoutes configures all HTTP routes for the server
func (s *Server) SetupRoutes() {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/validate", s.handleValidate)
	http.HandleFunc("/batch", s.handleBatch)
	http.HandleFunc("/validators", s.handleValidators)
	http.HandleFunc("/health", s.handleHealth)
}

// GetHandler returns the HTTP handler for the server
func (s *Server) GetHandler() http.Handler {
	// Add CORS middleware if enabled
	var handler http.Handler = http.DefaultServeMux
	return handler
}
