package utils

import (
	"encoding/json"
	"strings"
)

// parseJSONField extracts values from JSON-like fields (genres, spoken_languages)
func ParseJSONField(jsonStr string, key string) []string {
	jsonStr = strings.ReplaceAll(jsonStr, `'`, `"`)

	var items []map[string]any
	var result []string

	// Attempt to parse JSON
	err := json.Unmarshal([]byte(jsonStr), &items)
	if err != nil {
		return []string{jsonStr} // Return raw value if parsing fails
	}

	// Extract "name" fields
	for _, item := range items {
		if name, ok := item[key].(string); ok {
			result = append(result, name)
		}
	}
	return result
}

func ParseData(filename string) ([]map[string]string, error) {
	rows, err := ReadCSVFile(filename)
	if err != nil {
		return nil, err
	}

	// Extract header indices
	header := rows[0]
	var parsedData []map[string]string

	// Process rows
	for _, row := range rows[1:] {
		if len(row) != len(header) { // Check for inconsistent row length
			continue
		}

		data := make(map[string]string)
		for i, field := range row {
			data[header[i]] = field
		}

		parsedData = append(parsedData, data)
	}

	return parsedData, nil
}
