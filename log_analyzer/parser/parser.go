package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

type LogEntry struct {
	Level         string         `json:"level" validate:"required" mod:"default=unknown"`
	Message       string         `json:"msg" validate:"required" mod:"default=No message"`
	DynamicFields map[string]any `json:"-"` // Store extra fields
}

var (
	validate = validator.New()
	mod      = modifiers.New()
)

func ProcessLogFile(logFile string) ([]LogEntry, error) {
	data, err := os.ReadFile(logFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v\n", logFile, err)
	}

	var logEntries []LogEntry
	if err := json.Unmarshal(data, &logEntries); err != nil {
		return nil, fmt.Errorf("Error parsing JSON in %s: %v\n", logFile, err)

	}

	for i := range logEntries {
		if err := mod.Struct(nil, &logEntries[i]); err != nil {
			return nil, fmt.Errorf("error setting defaults: %v", err)
		}

		if err := validate.Struct(logEntries[i]); err != nil {
			return nil, fmt.Errorf("validation error: %v", err)
		}
	}

	return logEntries, nil
}
