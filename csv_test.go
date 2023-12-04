package logger_test

import (
	"testing"

	"github.com/threatwinds/logger"
)

func TestWriteCsvLog(t *testing.T) {
	headers := []string{"header1", "header2", "header3"}
	config := &logger.CSVConfig{Output: "file.csv", CsvHeaders: headers}
	logger, err := logger.NewCSVLogger(config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	values := map[string]string{
		"header1": "value1",
		"header2": "value2",
		"header3": "value3",
	}
	err = logger.WriteCsvLog(values)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify that the log entry is written to the CSV file
	// You can add additional assertions here based on your requirements
}
