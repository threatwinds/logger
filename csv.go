package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type CSVLogger struct {
	cnf *CSVConfig
}

type CSVConfig struct {
	Output     string      // <filepath>
	CsvHeaders []string    // csv fields
	csvFile    *os.File    // csv file
	csvWriter  *csv.Writer // csv writer
}

// WriteCsv writes a CSV record to the CSV file.
func NewCSVLogger(config *CSVConfig) (*CSVLogger, error) {
	var csvlogger = new(CSVLogger)
	if config == nil || config.CsvHeaders == nil || config.Output == "" {
		return nil, fmt.Errorf("config must be set")
	}

	if _, err := os.Stat(config.Output); !os.IsNotExist(err) {
		return nil, fmt.Errorf("file already exists")
	}

	// Create all directories in the path if they do not exist
	err := os.MkdirAll(filepath.Dir(config.Output), 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directories: %v", err)
	}

	csvFile, err := os.OpenFile(config.Output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	config.csvFile = csvFile
	config.csvWriter = csv.NewWriter(csvFile)
	defer config.csvWriter.Flush()

	err = config.csvWriter.Write(config.CsvHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to write to file: %v", err)
	}

	csvlogger.cnf = config
	return csvlogger, nil
}

// Close closes the CSV file.
func (l CSVLogger) Close() error {
	if l.cnf.csvFile != nil {
		return l.cnf.csvFile.Close()
	}
	return nil
}

// WriteCsvLog writes a log entry to a CSV file.
func (l CSVLogger) WriteCsvLog(values map[string]string) error {
	if l.cnf.csvFile != nil {
		var row []string
		for _, header := range l.cnf.CsvHeaders {
			if _, ok := values[header]; !ok {
				row = append(row, "")
			} else {
				row = append(row, values[header])
			}
		}

		l.cnf.csvWriter.Write(row)
		l.cnf.csvWriter.Flush()
	}
	return nil
}
