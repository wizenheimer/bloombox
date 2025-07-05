package main

import (
	"net/http"
	"os"

	"github.com/wizenheimer/bloombox/pkg/emailchecker"
	"github.com/wizenheimer/bloombox/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	defer logger.Sync()
	config := loadConfigFromEnv()

	checker, err := emailchecker.New(config)
	if err != nil {
		logger.Fatal("Failed to initialize email checker", zap.Error(err))
	}

	server := NewServer(checker, config)
	server.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Email validation API server starting", zap.String("port", port))
	logger.Info("Available endpoints:")
	logger.Info("  GET  /                 - API documentation")
	logger.Info("  POST /validate         - Validate single email")
	logger.Info("  POST /batch            - Validate multiple emails")
	logger.Info("  GET  /validators       - List available validators")
	logger.Info("  PUT  /validators/:name - Enable/disable validator")
	logger.Info("  GET  /health           - Health check")

	handler := server.GetHandler()
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
