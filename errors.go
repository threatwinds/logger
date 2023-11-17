package logger

import "strings"

// Error represents an error with a UUID, status, and message.
type Error struct {
	UUID    string `json:"uuid"`
	Status  int    `json:"-"`
	Message string `json:"error"`
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
