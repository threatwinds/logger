package rerror

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type Log struct {
	UUID     string `json:"uuid"`
	Code     int    `json:"code"`
	Severity string `json:"severity"`
	Path     string `json:"path"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
}

type Error struct {
	UUID    string
	HTTP    int
	RPC     codes.Code
	Message string
}

func LogF(code int, format string, args ...any) *Log {
	var log = new(Log)

	var severity string

	if code < 400 {
		severity = "INFO"
	} else if code >= 400 && code < 500 {
		severity = "WARNING"
	} else if code >= 500 && code < 600 {
		severity = "ERROR"
	} else {
		severity = "CRITICAL"
	}

	_, path, line, _ := runtime.Caller(2)

	log.UUID = uuid.NewString()
	log.Code = code
	log.Path = path
	log.Line = line
	log.Message = fmt.Sprintf(format, args...)
	log.Severity = severity

	j, _ := json.Marshal(log)
	fmt.Println(string(j))

	return log
}

func ErrorF(http int, rpc codes.Code, format string, args ...any) *Error {
	var l = LogF(http, format, args...)

	if l.Code >= 500 {
		return &Error{UUID: l.UUID, HTTP: http, RPC: rpc, Message: "internal"}
	}

	return &Error{UUID: l.UUID, HTTP: http, RPC: rpc, Message: l.Message}
}
