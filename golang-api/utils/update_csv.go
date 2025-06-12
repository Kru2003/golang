package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func SaveToCSV(fileName string, updatedRows [][]string) error {
	// Get the current working directory
	workingDirPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	// Create the file path by joining the working directory with the file name
	filePath := filepath.Join(workingDirPath, fileName)

	// Update the CSV with the new rows
	err = UpdateCSV(filePath, updatedRows)
	if err != nil {
		return fmt.Errorf("error updating %s: %v", fileName, err)
	}

	return nil
}

func UpdateCSV(filepath string, data [][]string) error {
	file, err := os.Create(filepath) // Overwrite the existing file
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure data is written to the file

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
