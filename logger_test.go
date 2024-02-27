package logger_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/threatwinds/logger"
)

func TestStdOut(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400, Output: "stdout"}

	logger := logger.NewLogger(config)

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

func TestJsonLogFile(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 100, Output: "test.json"}
	logger := logger.NewLogger(config)

	logger.LogF(200, "Test log message")

	// Check if file exists
	if _, err := os.Stat("test.json"); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// Check file content
	content, err := os.ReadFile("test.json")
	if err != nil {
		t.Error("Failed to read file")
	}

	if !strings.Contains(string(content), "Test log message") {
		t.Error("Log message was not written to file")
	}

	// Clean up
	os.Remove("test.json")
}

func TestCSVLogFile(t *testing.T) {
	config := &logger.Config{Format: "csv", Level: 100, Output: "test.log"}
	logger := logger.NewLogger(config)

	logger.LogF(200, "Test log message")

	// Check if file exists
	if _, err := os.Stat("test.log"); os.IsNotExist(err) {
		t.Error("File was not created")
	}

	// Check file content
	content, err := os.ReadFile("test.log")
	if err != nil {
		t.Error("Failed to read file")
	}

	if !strings.Contains(string(content), "Test log message") {
		t.Error("Log message was not written to file")
	}

	// Clean up
	os.Remove("test.log")
}

func TestInfo(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400}
	logger := logger.NewLogger(config)

	log := logger.Info("Test info message")

	if log == nil {
		t.Error("Log should not be nil")
	} else if log.Severity != "INFO" {
		t.Errorf("Expected severity to be INFO, got %s", log.Severity)
	}

	if log.Status != 200 {
		t.Errorf("Expected status code to be 200, got %d", log.Status)
	}
}

func TestRunWithRetries(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400}
	logger := logger.NewLogger(config)

	var result int

	err := logger.RunWithRetries(map[int][]string{http.StatusBadGateway: {"hella", "hello"}, http.StatusAccepted: {"helle"}}, func() error {
		result += 3
		return fmt.Errorf("hello")
	})

	if err == nil {
		t.Error("err should not be nil")
	}

	if result != 15 {
		t.Errorf("Expected 15, got %d", result)
	}
}
