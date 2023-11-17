package logger

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

// defaultConfig returns the default configuration for the logger.
func defaultConfig() *Config {
	return &Config{Format: "json", Level: 400}
}

// Log represents a log entry.
type Log struct {
	Timestamp string `json:"timestamp"`
	Severity  string `json:"severity"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Error
}

// Logger represents a logger instance.
type Logger struct{
	cnf *Config
}

// Config represents the configuration for the logger.
type Config struct {
	Format string // json, text
	Level int // 100: DEBUG, 200: INFO, 300: NOTICE, 400: WARNING, 500: ERROR, 502: CRITICAL, 509: ALERT
}

// New creates a new logger instance with the given configuration.
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

// ToJSON converts the log entry to JSON formated string.
func (l Log) ToJSON() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// ToString converts the log entry to a string format.
func (l Log) ToString() string {
	return fmt.Sprintf("%s %s %s %s %d %s", l.Timestamp, l.Severity, l.UUID, l.Path, l.Line, l.Message)
}

// LogF logs a formatted message with the given status code and arguments.
// It returns the log entry.
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
	} else if statusCode >= 502 && statusCode < 509 {
		severity = "CRITICAL"
	} else if statusCode >= 509 && statusCode < 511 {
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

// ErrorF logs an error message with the given status code and arguments.
// It returns an error object.
func (l Logger) ErrorF(statusCode int, format string, args ...any) *Error {
	var log = l.LogF(statusCode, format, args...)

	if log.Status >= 500 {
		return &Error{UUID: log.UUID, Status: statusCode, Message: "internal error"}
	}

	return &Error{UUID: log.UUID, Status: statusCode, Message: log.Message}
}
