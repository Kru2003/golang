package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func ReadCSVFile(filename string) ([][]string, error) {
	workingDirPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting working directory: %v", err)
	}

	// Fetch file path from config
	filePath := filepath.Join(workingDirPath, filename)

	// Open CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read CSV
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable column count to prevent errors

	// Read all rows
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("no data found in CSV")
	}

	return rows, nil
}
