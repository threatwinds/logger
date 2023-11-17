package rerror

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

func defaultConfig() *Config {
	return &Config{Format: "json", Level: 400}
}

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
}

type Logger struct{
	cnf *Config
}

type Config struct {
	Format string // json, text
	Level int // 100: DEBUG, 200: INFO, 300: NOTICE, 400: WARNING, 500: ERROR, 502: CRITICAL, 505: ALERT
}

func New(config *Config) *Logger{
	var logger = new(Logger)
	if config != nil {
		if config.Format == "" {
			config.Format = defaultConfig().Format
		}

		if config.Level == 0 {
			config.Level = defaultConfig().Level
		}
	} else {
		config = defaultConfig()
	}

	logger.cnf = config
	
	return logger
}

func (l Log) ToJSON() string {
	b, _ := json.Marshal(l)
	return string(b)
}

func (l Log) ToString() string {
	return fmt.Sprintf("%s %s %s %s %d %s", l.Timestamp, l.Severity, l.UUID, l.Path, l.Line, l.Message)
}

func (l Logger) LogF(statusCode int, format string, args ...any) *Log {
	var log = new(Log)

	var severity string

	if statusCode >= 100 && statusCode < 200 {
		severity = "DEBUG"
	} else if statusCode >= 200 && statusCode < 300 {
		severity = "INFO"
	} else if statusCode >= 300 && statusCode < 400 {
		severity = "NOTICE"
	} else if statusCode >= 400 && statusCode < 500 {
		severity = "WARNING"
	} else if statusCode >= 500 && statusCode < 502 {
		severity = "ERROR"
	} else if statusCode >= 502 && statusCode < 505 {
		severity = "CRITICAL"
	} else if statusCode >= 505 && statusCode < 510 {
		severity = "ALERT"
	} else {
		severity = "DEFAULT"
	}

	_, path, line, _ := runtime.Caller(2)

	log.UUID = uuid.NewString()
	log.Status = statusCode
	log.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	log.Path = path
	log.Line = line
	log.Message = fmt.Sprintf(format, args...)
	log.Severity = severity

	if statusCode >=  l.cnf.Level{
		var final string
		switch l.cnf.Format {
			case "json":
				final = log.ToJSON()
			case "text":{
				final = log.ToString()
			}
		}

		fmt.Println(final)
	}

	return log
}

func (l Logger) ErrorF(statusCode int, format string, args ...any) *Error {
	var log = l.LogF(statusCode, format, args...)

	if log.Status >= 500 {
		return &Error{UUID: log.UUID, Status: statusCode, Message: "internal error"}
	}

	return &Error{UUID: log.UUID, Status: statusCode, Message: log.Message}
}
