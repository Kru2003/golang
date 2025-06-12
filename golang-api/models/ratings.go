package models

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
)

type Ratings struct {
	UserId  string `json:"userId" validate:"required,gte=1"`
	MovieId string `json:"movieId" validate:"required"`
	Rating  string `json:"rating" validate:"required,gte=0,lte=5"`
}

type MovieRatings struct {
	MovieId string
	Ratings float64
}

type RatingModel struct {
	Ratings []Ratings
	loaded  bool
	mu      sync.Mutex
}

func NewRatingsModel() *RatingModel {
	return &RatingModel{}
}

// Function to load movies using utils.ParseData()
func (r *RatingModel) LoadRatings() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	config := config.AllConfig
	data, err := utils.ParseData(config.Ratings)
	if err != nil {
		return err
	}

	var ratings []Ratings
	for _, row := range data {
		ratings = append(ratings, Ratings{
			UserId:  row["userId"],
			MovieId: row["movieId"],
			Rating:  row["rating"],
		})
	}
	r.Ratings = ratings
	return nil

}

func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// Function to calculate ratings of a movie
func (r *RatingModel) CalculateAverageRatings() []MovieRatings {
	ratingSum := make(map[string]float64)
	ratingCount := make(map[string]int)

	// Aggregate ratings
	for _, rating := range r.Ratings {
		ratingValue, _ := strconv.ParseFloat(rating.Rating, 64)
		ratingSum[rating.MovieId] += ratingValue
		ratingCount[rating.MovieId]++
	}

	var movieRatings []MovieRatings
	for movieID, sum := range ratingSum {
		avgRating := RoundToTwoDecimals(sum / float64(ratingCount[movieID]))
		movieRatings = append(movieRatings, MovieRatings{
			MovieId: movieID,
			Ratings: avgRating,
		})
	}
	return movieRatings
}

// Function to list all movies ratings
func (r *RatingModel) ListRatings(page, limit int) ([]MovieRatings, error) {
	if !r.loaded {
		if err := r.LoadRatings(); err != nil {
			return nil, err
		}
		r.loaded = true
	}

	if len(r.Ratings) == 0 {
		return nil, errors.New("no ratings loaded, call LoadRatings first")
	}

	movieRatings := r.CalculateAverageRatings()

	return utils.Paginate(movieRatings, page, limit)
}

// Function to get ratings of a movie having a id movieId
func (r *RatingModel) GetRatingsByMovieId(movieId string) (MovieRatings, error) {
	if !r.loaded {
		if err := r.LoadRatings(); err != nil {
			return MovieRatings{}, err
		}
		r.loaded = true
	}

	if len(r.Ratings) == 0 {
		return MovieRatings{}, errors.New("no ratings loaded, call LoadRatings first")
	}

	movieRatings := r.CalculateAverageRatings()
	for _, rating := range movieRatings {
		if rating.MovieId == movieId {
			return rating, nil
		}
	}
	return MovieRatings{}, errors.New("movie not found")
}

// Function to add ratings
func (r *RatingModel) AddRatings(rating *Ratings) error {
	if !r.loaded {
		if err := r.LoadRatings(); err != nil {
			return err
		}
		r.loaded = true
	}

	r.Ratings = append(r.Ratings, *rating)

	workingDirPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %v", err)
	}

	// Fetch file path from config
	filePath := filepath.Join(workingDirPath, config.AllConfig.Ratings)

	err = appendRatingsToCSV(filePath, *rating)
	if err != nil {
		return fmt.Errorf("error in writing record: %v", err)
	}
	return nil
}

// Function to add ratings at the end of csv file
func appendRatingsToCSV(filepath string, rating Ratings) error {
	// Open the file in append mode
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Failed to open file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	timestamp := time.Now().Unix()
	// Convert struct to a slice of strings
	record := []string{
		rating.UserId,
		rating.MovieId,
		rating.Rating,
		fmt.Sprintf("%d", timestamp),
	}

	// Write the record to CSV
	if err := writer.Write(record); err != nil {
		return errors.New("Failed to write record")
	}

	log.Println("Rating appended successfully!")
	return nil
}

// Function to delete ratings of a movie
func (r *RatingModel) DeleteRatings(movieId string, userId *string) error {
	return r.ModifyRatings(userId, &movieId, nil, nil, "delete")
}

// UpdateRatings is to update the ratings of a movie
func (r *RatingModel) UpdateRatings(userId, movieId, newRating, newTimestamp string) error {
	return r.ModifyRatings(&userId, &movieId, &newRating, &newTimestamp, "update")
}

func (r *RatingModel) ModifyRatings(userId, movieId, newRating, newTimestamp *string, operation string) error {
	if !r.loaded {
		if err := r.LoadRatings(); err != nil {
			return err
		}
		r.loaded = true
	}

	rows, err := utils.ReadCSVFile(config.AllConfig.Ratings)
	if err != nil {
		return err
	}

	var updatedRows [][]string
	var updatedRatings []Ratings

	header := rows[0]
	updatedRows = append(updatedRows, header)
	modified := false

	for _, row := range rows[1:] {
		if len(row) < 4 {
			fmt.Println("Skipping invalid row:", row)
			continue
		}

		if row[1] == *movieId && (userId == nil || row[0] == *userId) {
			if operation == "delete" {
				modified = true
				continue
			} else if operation == "update" && newRating != nil && newTimestamp != nil {
				row[2] = *newRating
				row[3] = *newTimestamp
				modified = true
			}
		}

		updatedRows = append(updatedRows, row)
		updatedRatings = append(updatedRatings, Ratings{
			UserId:  row[0],
			MovieId: row[1],
			Rating:  row[2],
		})
	}

	if !modified {
		return errors.New("rating not found ")
	}

	r.Ratings = updatedRatings

	err = utils.SaveToCSV(config.AllConfig.Ratings, updatedRows)
	if err != nil {
		return fmt.Errorf("error updating Ratings.csv: %v", err)
	}

	return nil
}
