package logger

import (
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	config := zap.NewProductionConfig()

	// Set log level based on environment
	if os.Getenv("DEBUG") == "true" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configure encoding
	config.Encoding = "json"
	if os.Getenv("LOG_FORMAT") == "console" {
		config.Encoding = "console"
	}

	// Configure output paths
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

func Sync() {
	Logger.Sync()
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}
