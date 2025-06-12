package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
)

type CrewMember struct {
	CreditID   string `json:"credit_id"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Job        string `json:"job"`
}

// CrewModel handles all crew-related operations.
type CrewModel struct {
	CrewData map[int][]CrewMember
	loaded   bool
}

// NewCrewModel initializes a new CrewModel.
func NewCrewModel() *CrewModel {
	return &CrewModel{
		CrewData: make(map[int][]CrewMember),
	}
}

// ParseCrewData parses the crew data from credits CSV file.
func ParseCrewData() (map[int][]CrewMember, error) {
	records, err := utils.ReadCSVFile(config.AllConfig.Credits)
	if err != nil {
		return nil, err
	}

	movieCrews := make(map[int][]CrewMember)

	// Iterate over rows, skipping the header
	for _, row := range records[1:] {
		if len(row) < 3 {
			continue
		}

		var crew []CrewMember
		movieID, err := strconv.Atoi(strings.TrimSpace(row[2])) // Convert movieID from string to int
		if err != nil {
			fmt.Printf("Warning: Invalid movie ID '%s': %v\n", row[2], err)
			continue
		}

		row[1] = strings.ReplaceAll(row[1], `'`, `"`)

		err = json.Unmarshal([]byte(row[1]), &crew)
		if err != nil {
			fmt.Printf("Warning: Error parsing crew JSON for movie %d: %v\n", movieID, err)
			continue
		}

		movieCrews[movieID] = crew
	}

	return movieCrews, nil

}

// LoadCrew loads crew data from the CSV and stores it in CrewModel.
func (c *CrewModel) LoadCrew() error {
	movieCrews, err := ParseCrewData()
	if err != nil {
		return err
	}

	c.CrewData = movieCrews
	c.loaded = true

	return nil
}

// ListCrewMembers is to list crew members of movie having id movieID
func (c *CrewModel) ListCrewMembers(movieID string) ([]CrewMember, error) {
	if !c.loaded {
		err := c.LoadCrew()
		if err != nil {
			return nil, err
		}
	}

	id, err := strconv.Atoi(movieID) // Convert string to int
	if err != nil {
		return nil, errors.New("invalid movie ID format")
	}

	movieCrew, exists := c.CrewData[id]
	if !exists {
		return nil, errors.New("movie not found")
	}

	return movieCrew, nil
}

// DeleteCreditsForMovie is to delete credits of a movie when that movie is deleted
func DeleteCreditsForMovie(movieId string) error {
	rows, err := utils.ReadCSVFile(config.AllConfig.Credits)
	if err != nil {
		return err
	}

	var updatedRows [][]string
	header := rows[0]
	updatedRows = append(updatedRows, header)
	deleted := false

	for _, row := range rows[1:] {
		if len(row) < 3 {
			fmt.Println("Skipping row due to insufficient columns:", row)
			continue
		}

		if row[2] == movieId {
			deleted = true
			continue
		}
		updatedRows = append(updatedRows, row)
	}

	if !deleted {
		fmt.Println("No credits found for movie:", movieId)
	}

	err = utils.SaveToCSV(config.AllConfig.Credits, updatedRows)
	if err != nil {
		return fmt.Errorf("error updating credits.csv: %v", err)
	}

	return nil
}

// UpdateCrewMember is to update details of crew member having crewId in a movie having movieId
func (c *CrewModel) UpdateCrewMember(movieId, crewId string, updatedCrew CrewMember) error {
	rows, err := utils.ReadCSVFile(config.AllConfig.Credits)
	if err != nil {
		return err
	}

	var updatedRows [][]string
	header := rows[0]
	updatedRows = append(updatedRows, header)
	updated := false
	crewDataMap := make(map[int][]CrewMember)

	for _, row := range rows[1:] {
		if len(row) < 3 {
			fmt.Println("Skipping row due to insufficient columns:", row)
			continue
		}

		// Parse movie ID as int
		movieIDInt, err := strconv.Atoi(row[2])
		if err != nil {
			fmt.Printf("Skipping row due to invalid movie ID: %s\n", row[2])
			continue
		}

		var crew []map[string]any                     // Use map to keep all fields intact
		row[1] = strings.ReplaceAll(row[1], `'`, `"`) // Ensure valid JSON format

		// Parse JSON into []map[string]any (preserving extra fields)
		err = json.Unmarshal([]byte(row[1]), &crew)
		if err != nil {
			fmt.Printf("Warning: Error parsing crew JSON for movie %s: %v\n", row[2], err)
			continue
		}

		// Find and update the crew member while keeping extra fields
		for i, member := range crew {
			idFloat, ok := member["id"].(float64) // JSON unmarshals numbers as float64
			if !ok {
				continue
			}
			if strconv.Itoa(int(idFloat)) == crewId { // Compare as string
				if updatedCrew.Name != "" {
					member["name"] = updatedCrew.Name
				}
				if updatedCrew.Department != "" {
					member["department"] = updatedCrew.Department
				}
				if updatedCrew.Job != "" {
					member["job"] = updatedCrew.Job
				}
				crew[i] = member // Update only relevant fields
				updated = true
				break
			}
		}

		// Convert back to JSON string
		crewJSON, err := json.Marshal(crew)
		if err != nil {
			return fmt.Errorf("error marshaling updated crew data: %v", err)
		}

		row[1] = string(crewJSON) // Store updated crew JSON in CSV
		updatedRows = append(updatedRows, row)

		// Convert []map[string]any back to []CrewMember (while keeping extra fields)
		var crewMembers []CrewMember
		for _, member := range crew {
			memberJSON, _ := json.Marshal(member) // Convert map to JSON
			var crewMember CrewMember
			json.Unmarshal(memberJSON, &crewMember) // Convert JSON to struct
			crewMembers = append(crewMembers, crewMember)
		}

		// Store in CrewData map
		crewDataMap[movieIDInt] = crewMembers
	}

	if !updated {
		return errors.New("crew member not found for the given movie ID")
	}

	c.CrewData = crewDataMap

	err = utils.SaveToCSV(config.AllConfig.Credits, updatedRows)
	if err != nil {
		return fmt.Errorf("error updating credits.csv: %v", err)
	}

	return nil
}
