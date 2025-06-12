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

type CastMember struct {
	CreditID  string `json:"credit_id"`
	ID        int    `json:"id"`
	Character string `json:"character"`
	Name      string `json:"name"`
}

type CastModel struct {
	CastData map[int][]CastMember
	loaded   bool
}

func NewCastModel() *CastModel {
	return &CastModel{
		CastData: make(map[int][]CastMember),
	}
}

func ParseCastsData() (map[int][]CastMember, error) {
	records, err := utils.ReadCSVFile(config.AllConfig.Credits)
	if err != nil {
		return nil, err
	}

	movieCasts := make(map[int][]CastMember)

	// Iterate over rows, skipping the header
	for _, row := range records[1:] {
		if len(row) < 3 {
			continue
		}

		var cast []CastMember
		movieID, err := strconv.Atoi(strings.TrimSpace(row[2])) // Convert movieID from string to int
		if err != nil {
			fmt.Printf("Warning: Invalid movie ID '%s': %v\n", row[2], err)
			continue
		}

		row[0] = strings.ReplaceAll(row[0], `'`, `"`)

		err = json.Unmarshal([]byte(row[0]), &cast)
		if err != nil {
			fmt.Printf("Warning: Error parsing cast JSON for movie %d: %v\n", movieID, err)
			continue
		}

		movieCasts[movieID] = cast
	}

	return movieCasts, nil
}

func (c *CastModel) LoadCast() error {
	movieCasts, err := ParseCastsData()
	if err != nil {
		return err
	}

	c.CastData = movieCasts
	c.loaded = true

	return nil
}

func (c *CastModel) ListCastMembers(movieID string) ([]CastMember, error) {
	if !c.loaded {
		err := c.LoadCast()
		if err != nil {
			return nil, err
		}
	}

	id, err := strconv.Atoi(movieID) // Convert string to int
	if err != nil {
		return nil, errors.New("invalid movie ID format")
	}

	movieCast, exists := c.CastData[id]
	if !exists {
		return nil, errors.New("movie not found")
	}
	return movieCast, nil
}

func (c *CastModel) ListMoviesByCastId(castId string) ([]int, error) {
	if !c.loaded {
		err := c.LoadCast()
		if err != nil {
			return nil, err
		}
	}

	id, err := strconv.Atoi(castId) // Convert string to int
	if err != nil {
		return nil, errors.New("invalid cast ID format")
	}

	var movieIDs []int

	// Iterate over all movies to find occurrences of the given cast member
	for movieID, castMembers := range c.CastData {
		for _, member := range castMembers {
			if member.ID == id {
				movieIDs = append(movieIDs, movieID)
				break // Avoid duplicate entries for the same movie
			}
		}
	}

	if len(movieIDs) == 0 {
		return nil, errors.New("no movies found for the given cast ID")
	}

	return movieIDs, nil
}

// UpdateCastMember is to update cast member details having castId of a movie having movieId
func (c *CastModel) UpdateCastMember(movieId, castId string, updatedCast CastMember) error {
	rows, err := utils.ReadCSVFile(config.AllConfig.Credits)
	if err != nil {
		return err
	}

	var updatedRows [][]string
	header := rows[0]
	updatedRows = append(updatedRows, header)
	updated := false
	castDataMap := make(map[int][]CastMember)

	for _, row := range rows[1:] {
		if len(row) < 3 {
			fmt.Println("Skipping row due to insufficient columns:", row)
			continue
		}

		movieIDInt, err := strconv.Atoi(row[2])
		if err != nil {
			fmt.Printf("Skipping row due to invalid movie ID: %s\n", row[2])
			continue
		}

		var cast []map[string]any                     // Use map to preserve extra fields
		row[0] = strings.ReplaceAll(row[0], `'`, `"`) // Ensure valid JSON format

		err = json.Unmarshal([]byte(row[0]), &cast)
		if err != nil {
			fmt.Printf("Warning: Error parsing crew JSON for movie %s: %v\n", movieId, err)
			continue
		}

		for i, member := range cast {
			if strconv.Itoa(int(member["id"].(float64))) == castId {
				if updatedCast.Name != "" {
					member["name"] = updatedCast.Name
				}
				if updatedCast.Character != "" {
					member["character"] = updatedCast.Character
				}

				cast[i] = member
				updated = true
				break
			}
		}

		castJSON, err := json.Marshal(cast)
		if err != nil {
			return fmt.Errorf("error marshaling updated crew data: %v", err)
		}

		row[1] = string(castJSON)
		updatedRows = append(updatedRows, row)

		var castMembers []CastMember
		for _, member := range cast {
			memberJSON, _ := json.Marshal(member) // Convert map to JSON
			var castMember CastMember
			json.Unmarshal(memberJSON, &castMember) // Convert JSON to struct
			castMembers = append(castMembers, castMember)
		}

		castDataMap[movieIDInt] = castMembers
	}

	if !updated {
		return errors.New("cast member not found for the given movie ID")
	}

	c.CastData = castDataMap

	err = utils.SaveToCSV(config.AllConfig.Credits, updatedRows)
	if err != nil {
		return fmt.Errorf("error updating credits.csv: %v", err)
	}

	return nil
}
