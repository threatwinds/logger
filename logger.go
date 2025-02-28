package logger

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/natefinch/lumberjack.v2"
)

// defaultConfig returns a pointer to the default Config instance with predefined values.
func defaultConfig() *Config {
	return &Config{
		Format:    "json",
		Level:     400,
		Output:    "stdout",
		Retries:   5,
		Wait:      1 * time.Second,
		StatusMap: map[int][]string{},
	}
}

// Log represents a structured log entry containing timestamp, severity, file path, line number, and error details.
type Log struct {
	Timestamp string `json:"timestamp"`
	Severity  string `json:"severity"`
	Path      string `json:"path"`
	Line      int    `json:"line"`
	Error
}

// Logger provides logging functionalities with configurable format, level, output, retries, and wait time.
type Logger struct {
	cnf *Config
}

// Config defines the configuration options for logging, including format, level, output, retries, and wait duration.
type Config struct {
	Format    string           // json, text, csv
	Level     int              // 100: DEBUG, 200: INFO, 300: NOTICE, 400: WARNING, 500: ERROR, 502: CRITICAL, 509: ALERT
	Output    string           // stdout, <filepath>
	Retries   int              // number of retries
	Wait      time.Duration    // wait time between retries
	StatusMap map[int][]string // status code to message map
}

// NewLogger initializes and returns a new Logger instance configured with the provided Config or default values.
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

		if config.Retries == 0 {
			config.Retries = defaultConfig().Retries
		}

		if config.Wait == 0 {
			config.Wait = defaultConfig().Wait
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

// ToJSON serializes the Log struct into a JSON-formatted string and returns it.
func (l *Log) ToJSON() string {
	b, _ := json.Marshal(l)
	return string(b)
}

// ToString converts the Log struct into a single-line string with fields concatenated and separated by spaces.
func (l *Log) ToString() string {
	return strings.Join([]string{l.Timestamp, l.Severity, l.UUID, l.Path, fmt.Sprint(l.Line), l.Message}, " ")
}

// ToCsv converts a Log instance into a CSV-formatted single-line string and returns it.
// Returns an empty string on failure.
func (l *Log) ToCsv() string {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	err := w.Write([]string{l.Timestamp, l.Severity, l.UUID, l.Path, fmt.Sprint(l.Line), l.Message})
	if err != nil {
		return ""
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return ""
	}

	return strings.TrimRight(b.String(), "\n")
}

// LogF logs a message with a specific status code, format string, and arguments, returning a structured Log instance.
// Determines severity based on the status code and outputs in the configured format and destination.
func (l *Logger) LogF(statusCode int, format string, args ...any) *Log {
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

// Info logs an informational message with a status code of 200, using formatted strings and optional arguments.
func (l *Logger) Info(format string, args ...any) *Log {
	return l.LogF(200, format, args...)
}
