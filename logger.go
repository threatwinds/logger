package logger

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/natefinch/lumberjack.v2"
)

// defaultConfig returns the default configuration for the logger.
func defaultConfig() *Config {
	return &Config{Format: "json", Level: 400, Output: "stdout"}
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
type Logger struct {
	cnf *Config
}

// Config represents the configuration for the logger.
type Config struct {
	Format string // json, text, csv
	Level  int    // 100: DEBUG, 200: INFO, 300: NOTICE, 400: WARNING, 500: ERROR, 502: CRITICAL, 509: ALERT
	Output string // stdout, <filepath>
}

// New creates a new logger instance with the given configuration.
func NewLogger(config *Config) *Logger {
	var logger = new(Logger)
	if config != nil {
		if config.Format == "" {
			config.Format = defaultConfig().Format
		}

		if config.Level == 0 {
			config.Level = defaultConfig().Level
		}

		if config.Output == "" {
			config.Output = defaultConfig().Output
		}
	} else {
		config = defaultConfig()
	}

	logger.cnf = config

	if logger.cnf.Output != "stdout" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logger.cnf.Output,
			MaxSize:    5, // megabytes
			MaxBackups: 100,
			MaxAge:     30, // days
		})

		log.SetFlags(0)
	}

	return logger
}

// ToJSON converts the log entry to JSON formated string.
func (l Log) ToJSON() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// ToString converts the log entry to a string format.
func (l Log) ToString() string {
	return strings.Join([]string{l.Timestamp, l.Severity, l.UUID, l.Path, fmt.Sprint(l.Line), l.Message}, " ")
}

// ToCsv converts the log entry to a CSV format.
func (l Log) ToCsv() string {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	err := w.Write([]string{l.Timestamp, l.Severity, l.UUID, l.Path, fmt.Sprint(l.Line), l.Message})
	if err != nil {
		return ""
	}

	w.Flush()

	if err := w.Error(); err != nil{
		return ""
	}

	return strings.TrimRight(b.String(), "\n")
}

// LogF logs a formatted message with the given status code and arguments.
// It returns the log entry.
func (l Logger) LogF(statusCode int, format string, args ...any) *Log {
	var newLog = new(Log)

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

	newLog.UUID = uuid.NewString()
	newLog.Status = statusCode
	newLog.Path = path
	newLog.Line = line
	newLog.Message = fmt.Sprintf(format, args...)
	newLog.Severity = severity
	newLog.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)

	if statusCode >= l.cnf.Level {
		var final string
		switch l.cnf.Format {
		case "json":
			final = newLog.ToJSON()
		case "text":
			final = newLog.ToString()
		case "csv":
			final = newLog.ToCsv()
		}

		switch l.cnf.Output {
		case "stdout":
			fmt.Println(final)
		default:
			log.Println(final)
		}
	}

	return newLog
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

// Fatal logs an error message and exits the program
func (l Logger) Fatal(format string, args ...any) {
	l.LogF(501, format, args...)
	os.Exit(1)
}

// Info logs an info message with the given arguments.
func (l Logger) Info(format string, args ...any) *Log {
	return l.LogF(200, format, args...)
}
