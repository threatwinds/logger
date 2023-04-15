package rerror

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type Log struct {
	Timestamp string `json:"timestamp"`
	UUID      string `json:"uuid"`
	Status    int    `json:"status"`
	Severity  string `json:"severity"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Message   string `json:"message"`
}

type Error struct {
	UUID    string
	Status  int
	Message string
}

func LogF(status int, format string, args ...any) *Log {
	var log = new(Log)

	var severity string

	if status < 400 {
		severity = "INFO"
	} else if status >= 400 && status < 500 {
		severity = "WARNING"
	} else if status >= 500 && status < 600 {
		severity = "ERROR"
	} else {
		severity = "CRITICAL"
	}

	_, path, line, _ := runtime.Caller(2)

	log.UUID = uuid.NewString()
	log.Status = status
	log.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	log.Path = path
	log.Line = line
	log.Message = fmt.Sprintf(format, args...)
	log.Severity = severity

	j, _ := json.Marshal(log)
	fmt.Println(string(j))

	return log
}

func ErrorF(status int, format string, args ...any) *Error {
	var l = LogF(status, format, args...)

	if l.Status >= 500 {
		return &Error{UUID: l.UUID, Status: status, Message: "internal error"}
	}

	return &Error{UUID: l.UUID, Status: status, Message: l.Message}
}
