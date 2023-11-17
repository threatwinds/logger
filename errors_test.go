package logger_test

import (
	"errors"
	"testing"

	"github.com/threatwinds/logger"
)

func TestIs(t *testing.T) {
	err := errors.New("test error")

	// Test case 1: Error contains the specified substring
	if !logger.Is(err, "test") {
		t.Error("Expected Is to return true")
	}

	// Test case 2: Error does not contain the specified substring
	if logger.Is(err, "foo") {
		t.Error("Expected Is to return false")
	}
}

func TestError_Is(t *testing.T) {
	err := &logger.Error{
		UUID:    "123",
		Status:  500,
		Message: "internal server error",
	}

	// Test case 1: Error message contains the specified substring
	if !err.Is("server") {
		t.Error("Expected Is to return true")
	}

	// Test case 2: Error message does not contain the specified substring
	if err.Is("foo") {
		t.Error("Expected Is to return false")
	}
}