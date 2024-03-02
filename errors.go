package logger

import (
	"fmt"
	"os"
	"strings"
)

// Error represents an error with a UUID, status, and message.
type Error struct {
	UUID    string `json:"uuid"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Is checks if the given error contains any of the specified substrings.
// It returns true if any of the substrings are found, otherwise false.
func Is(e error, kargs ...string) bool {
	for _, arg := range kargs {
		if strings.Contains(e.Error(), arg) {
			return true
		}
	}
	return false
}

// Is checks if the error's message contains any of the specified substrings.
// It returns true if any of the substrings are found, otherwise false.
func (e *Error) Is(kargs ...string) bool {
	for _, arg := range kargs {
		if strings.Contains(e.Message, arg) {
			return true
		}
	}
	return false
}

// ErrorF logs an error message with the given format and arguments.
// It returns an error object.
func (l Logger) ErrorF(format string, args ...any) *Error {
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

// Fatal logs an error message and exits the program
func (l Logger) Fatal(format string, args ...any) {
	l.LogF(500, format, args...)
	os.Exit(1)
}
