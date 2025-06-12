package models

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
)

type Movies struct {
	ID               string   `json:"id" validate:"required"`
	OriginalLanguage string   `json:"original_language" validate:"required,min=2"`
	Title            string   `json:"title" validate:"required,min=2,max=100"`
	Popularity       string   `json:"popularity" validate:"gte=0,lte=100"`
	Genres           []string `json:"genres" validate:"required,dive,min=5,max=50"`
	ReleaseDate      string   `json:"release_date" validate:"required,datetime=2006-01-02"`
	Runtime          string   `json:"runtime" validate:"gte=0,lte=100"`
	SpokenLanguages  []string `json:"spoken_languages" validate:"required,dive,min=5,max=50"`
	Status           string   `json:"status" validate:"required,oneof=Released Upcoming Cancelled"`
}

type MovieModel struct {
	Movies []Movies
	loaded bool
	mu     sync.Mutex
}

func NewMovieModel() *MovieModel {
	return &MovieModel{}
}

// Function to load movies using utils.ParseData()
func (m *MovieModel) LoadMovies() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	config := config.AllConfig
	data, err := utils.ParseData(config.Movies)
	if err != nil {
		return err
	}

	var movies []Movies
	for _, row := range data {
		movies = append(movies, Movies{
			ID:               row["id"],
			OriginalLanguage: row["original_language"],
			Title:            row["title"],
			Popularity:       row["popularity"],
			Genres:           utils.ParseJSONField(row["genres"], "name"),
			ReleaseDate:      row["release_date"],
			Runtime:          row["runtime"],
			SpokenLanguages:  utils.ParseJSONField(row["spoken_languages"], "name"),
			Status:           row["status"],
		})
	}
	m.Movies = movies
	return nil
}

// ListMovies fetches paginated movies
func (m *MovieModel) ListMovies(filters map[string]string, page, limit int) ([]Movies, error) {
	if !m.loaded {
		if err := m.LoadMovies(); err != nil {
			return nil, err
		}
		m.loaded = true
	}

	if len(m.Movies) == 0 {
		return nil, errors.New("no movies loaded, call LoadMovies first")
	}

	movieName := strings.ToLower(filters["name"])
	movieGenre := strings.ToLower(filters["genre"])
	movieLanguage := strings.ToLower(filters["language"])

	var matchedMovies []Movies

	for _, movie := range m.Movies {
		if movieName != "" && !strings.Contains(strings.ToLower(movie.Title), strings.ToLower(movieName)) {
			continue
		}

		if movieGenre != "" {
			found := false
			for _, genre := range movie.Genres {
				if strings.ToLower(genre) == strings.ToLower(movieGenre) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if movieLanguage != "" {
			found := false
			for _, language := range movie.SpokenLanguages {
				if strings.ToLower(language) == strings.ToLower(movieLanguage) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// If all conditions pass, add movie to matched list
		matchedMovies = append(matchedMovies, movie)

	}

	// If no movies match the filters, return an error
	if len(matchedMovies) == 0 {
		return nil, errors.New("no movies found matching the given criteria")
	}

	return utils.Paginate(matchedMovies, page, limit)
}

// function GetMovie to get movie by its specified movieID
func (m *MovieModel) GetMovie(movieID string) (Movies, error) {
	if !m.loaded {
		if err := m.LoadMovies(); err != nil {
			return Movies{}, err
		}
		m.loaded = true
	}
	for i := range m.Movies {
		if m.Movies[i].ID == movieID {
			return m.Movies[i], nil
		}
	}
	return Movies{}, errors.New("movie not found")
}

// Function is to add movie
func (m *MovieModel) AddMovie(movie *Movies) error {
	if !m.loaded {
		if err := m.LoadMovies(); err != nil {
			return err
		}
		m.loaded = true
	}
	for _, existingMovie := range m.Movies {
		if existingMovie.ID == movie.ID || existingMovie.Title == movie.Title {
			return errors.New("movie with this ID or title already exists")
		}
	}

	// Add new movie to the list
	m.Movies = append(m.Movies, *movie)
	workingDirPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	// Fetch file path from config
	filePath := filepath.Join(workingDirPath, config.AllConfig.Movies)

	err = appendMovieToCSV(filePath, *movie)
	if err != nil {
		return err
	}
	return nil
}

func formatData(datas []string) string {
	var formattedData []string
	for _, data := range datas {
		formattedData = append(formattedData, fmt.Sprintf("{'name': '%s'}", data))
	}
	return fmt.Sprintf("[%s]", joinStrings(formattedData, ", "))
}

func joinStrings(elements []string, sep string) string {
	result := ""
	for i, element := range elements {
		result += element
		if i < len(elements)-1 {
			result += sep
		}
	}
	return result
}

// appendMovieToCSV is to append movies at the end of csv file
func appendMovieToCSV(filePath string, movie Movies) error {
	// Open the file in append mode or create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Convert slice fields (Genres & SpokenLanguages) to string
	genreJSON := formatData(movie.Genres)
	spokenLanguagesStr := formatData(movie.SpokenLanguages)

	// Convert struct to a slice of strings
	record := []string{
		"",
		"",
		"",
		genreJSON,
		"",
		movie.ID,
		"",
		movie.OriginalLanguage,
		"",
		"",
		movie.Popularity,
		"",
		"",
		"",
		movie.ReleaseDate,
		"",
		movie.Runtime,
		spokenLanguagesStr,
		movie.Status,
		"",
		movie.Title,
		"",
		"",
		"",
	}

	// Write the record to CSV
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("Failed to write record: %v", err)
	}

	log.Println("Movie appended successfully!")
	return nil
}

func (m *MovieModel) MovieExists(movieId string) (bool, error) {
	if !m.loaded {
		if err := m.LoadMovies(); err != nil {
			return false, err
		}
		m.loaded = true
	}
	for _, movie := range m.Movies {
		if movie.ID == movieId {
			return true, nil
		}
	}
	return false, nil
}

func parseCSVArray(data string) []string {
	data = strings.Trim(data, "[]") // Remove square brackets
	if data == "" {
		return []string{}
	}
	parts := strings.Split(data, ",") // Split by comma
	for i := range parts {
		parts[i] = strings.Trim(parts[i], " '\"") // Trim spaces and quotes
	}
	return parts
}

// ModifyMovie will modify movies according to input in struct as well as in csv file
func (m *MovieModel) ModifyMovie(movieId string, updatedMovie *Movies, operation string) error {
	if !m.loaded {
		if err := m.LoadMovies(); err != nil {
			return err
		}
		m.loaded = true
	}

	rows, err := utils.ReadCSVFile(config.AllConfig.Movies)
	if err != nil {
		return err
	}

	var updatedRows [][]string
	var updatedMovies []Movies

	header := rows[0]
	updatedRows = append(updatedRows, header)
	modified := false

	for _, row := range rows[1:] {
		if len(row) < 23 {
			fmt.Println("Skipping invalid row:", row)
			continue
		}

		genres := parseCSVArray(row[3])
		spokenlanguages := parseCSVArray(row[17])

		if row[5] == movieId {
			if operation == "delete" {
				modified = true
				continue
			} else if operation == "update" && updatedMovie != nil {
				row[5] = movieId
				row[7] = updatedMovie.OriginalLanguage
				row[10] = updatedMovie.Popularity
				row[14] = updatedMovie.ReleaseDate
				row[16] = updatedMovie.Runtime
				row[18] = updatedMovie.Status
				row[20] = updatedMovie.Title
				row[3] = formatData(updatedMovie.Genres)
				row[17] = formatData(updatedMovie.SpokenLanguages)
				modified = true
			}
		}

		updatedRows = append(updatedRows, row)

		updatedMovies = append(updatedMovies, Movies{
			ID:               row[5],
			OriginalLanguage: row[7],
			Title:            row[20],
			Popularity:       row[10],
			Genres:           genres,
			ReleaseDate:      row[14],
			Runtime:          row[16],
			SpokenLanguages:  spokenlanguages,
			Status:           row[18],
		})
	}

	if !modified {
		return errors.New("movie not found")
	}

	m.Movies = updatedMovies

	err = utils.SaveToCSV(config.AllConfig.Movies, updatedRows)
	if err != nil {
		return fmt.Errorf("error updating movies_metadata.csv: %v", err)
	}

	if operation == "delete" {
		ratingModel := &RatingModel{}
		if err := ratingModel.DeleteRatings(movieId, nil); err != nil {
			return fmt.Errorf("error deleting ratings for movie: %v", err)
		}
		err = DeleteCreditsForMovie(movieId)
		if err != nil {
			return fmt.Errorf("error deleting credits: %v", err)
		}
	}

	return nil
}

// DeleteMovie is to delete movie having movieId
func (m *MovieModel) DeleteMovie(movieId string) error {
	return m.ModifyMovie(movieId, nil, "delete")
}

// UpdateMovie is to update movie having movieId
func (m *MovieModel) UpdateMovie(movieId string, updatedMovie *Movies) error {
	return m.ModifyMovie(movieId, updatedMovie, "update")
}
