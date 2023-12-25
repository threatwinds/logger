package logger_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/threatwinds/logger"
)

func TestLogF(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400}
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

func TestErrorF(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 400}
	logger := logger.NewLogger(config)

	L := logger.ErrorF(500, "Test error message")

	if L == nil {
		t.Error("Error should not be nil")
	} else if L.Status != 500 {
		t.Errorf("Expected status code to be 500, got %d", L.Status)
	}

	if L.Message != "internal error" {
		t.Errorf("Expected error message to be 'internal error', got '%s'", L.Message)
	}
}

func TestLogFWithOutput(t *testing.T) {
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

func TestJsonFileCreation(t *testing.T) {
	config := &logger.Config{Format: "json", Level: 100, Output: "test.json"}
	logger := logger.NewLogger(config)

	logger.LogF(200, "Test log message")

	// Check if file exists
	if _, L := os.Stat("test.json"); os.IsNotExist(L) {
		t.Error("File was not created")
	}

	// Check file content
	content, L := ioutil.ReadFile("test.json")
	if L != nil {
		t.Error("Failed to read file")
	}

	if !strings.Contains(string(content), "Test log message") {
		t.Error("Log message was not written to file")
	}

	// Clean up
	os.Remove("test.json")
}

func TestLogFileCreation(t *testing.T) {
	config := &logger.Config{Format: "text", Level: 100, Output: "test.log"}
	logger := logger.NewLogger(config)

	logger.LogF(200, "Test log message")

	// Check if file exists
	if _, L := os.Stat("test.log"); os.IsNotExist(L) {
		t.Error("File was not created")
	}

	// Check file content
	content, L := ioutil.ReadFile("test.log")
	if L != nil {
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
func TestCsvLog(t *testing.T) {
	filename := "test.csv"
	config := &logger.Config{Format: "csv", Level: 200, Output: filename}
	logger := logger.NewLogger(config)

	type TestData struct {
		Field1 string
		Field2 string
		Field3 string
	}

	data := TestData{
		Field1: "value1",
		Field2: "value2",
		Field3: "value3",
	}

	err := logger.CsvLog(&data)
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	// Check if CSV file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("CSV file was not created")
	}

	// Check file content
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Error("Failed to read CSV file")
	}

	expectedHeader := "Field1,Field2,Field3"
	if !strings.Contains(string(content), expectedHeader) {
		t.Errorf("Expected CSV file to contain header '%s'", expectedHeader)
	}

	expectedRow := "value1,value2,value3"
	if !strings.Contains(string(content), expectedRow) {
		t.Errorf("Expected CSV file to contain row '%s'", expectedRow)
	}

	// Clean up
	os.Remove(filename)
}
