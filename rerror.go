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
	Severity  string `json:"severity"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Error
}

type Error struct {
	UUID    string `json:"uuid"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
}

func LogF(print bool, status int, tag, format string, args ...any) *Log {
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
	log.Tag = tag

	if print {
		j, _ := json.Marshal(log)
		fmt.Println(string(j))
	}

	return log
}

func ErrorF(print bool, status int, tag, format string, args ...any) *Error {
	var l = LogF(print, status, tag, format, args...)

	if l.Status >= 500 {
		return &Error{UUID: l.UUID, Status: status, Message: "internal error", Tag: l.Tag}
	}

	return &Error{UUID: l.UUID, Status: status, Message: l.Message, Tag: l.Tag}
}
