package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Error represents an error with a unique identifier, status code, and a message.
type Error struct {
	UUID    string `json:"uuid"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error serializes the Error struct into a JSON-formatted string and returns it as an error message.
func (e *Error) Error() string {
	a, _ := json.Marshal(e)
	return string(a)
}

// Is checks if the provided error message contains any of the specified substrings and returns true if a match is found.
func Is(e error, args ...string) bool {
	for _, arg := range args {
		if strings.Contains(e.Error(), arg) {
			return true
		}
	}
	return false
}

// IsAny checks if the error message contains any of the provided substrings and returns true if a match is found.
func (e *Error) IsAny(args ...string) bool {
	for _, arg := range args {
		if strings.Contains(e.Message, arg) {
			return true
		}
	}
	return false
}

// ErrorF logs an error message with a formatted string and arguments and determines the corresponding status code.
// Returns an error instance containing a UUID, status code, and message.
func (l *Logger) ErrorF(format string, args ...any) *Error {
	var statusCode int = 500

	for k, v := range l.cnf.StatusMap {
		for _, msg := range v {
			if strings.Contains(fmt.Sprintf(format, args...), msg) {
				statusCode = k
				break
			}
		}
	}

	var log = l.LogF(statusCode, format, args...)

	if log.Status >= 500 {
		return &Error{UUID: log.UUID, Status: statusCode, Message: "internal error"}
	}

	return &Error{UUID: log.UUID, Status: statusCode, Message: log.Message}
}

// Fatal logs a message with a status code of 500 and terminates the application with os.Exit(1).
func (l *Logger) Fatal(format string, args ...any) {
	_ = l.LogF(500, format, args...)
	os.Exit(1)
}
