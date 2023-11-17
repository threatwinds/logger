package logger_test

import (
	"testing"
	"github.com/threatwinds/logger"
)

func TestLogF(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400}
	logger := logger.New(config)

	log := logger.LogF(200, "Test log message")

	if log == nil {
		t.Error("Log should not be nil")
	} else if log.Severity != "INFO" {
		t.Errorf("Expected severity to be INFO, got %s", log.Severity)
	}

	if log.Status != 200 {
		t.Errorf("Expected status code to be 200, got %d", log.Status)
	}
}

func TestErrorF(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400}
	logger := logger.New(config)

	err := logger.ErrorF(500, "Test error message")

	if err == nil {
		t.Error("Error should not be nil")
	} else if err.Status != 500 {
		t.Errorf("Expected status code to be 500, got %d", err.Status)
	}

	if err.Message != "internal error" {
		t.Errorf("Expected error message to be 'internal error', got '%s'", err.Message)
	}
}