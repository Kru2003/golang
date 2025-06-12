package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
)

// RatingsTable represent table name
const RatingsTable = "ratings"

type Ratings struct {
	UserId  int     `db:"user_id" json:"userId" validate:"required"`
	MovieId int     `db:"movie_id" json:"movieId" validate:"required"`
	Rating  float32 `db:"rating" json:"rating" validate:"required"`
}

type MovieRating struct {
	MovieId int     `db:"movie_id"`
	Title   string  `db:"title"`
	Rating  float32 `db:"avg_rating"`
}

type RatingModel struct {
	db *goqu.Database
}

func InitRatingsModel(goqu *goqu.Database) (*RatingModel, error) {
	return &RatingModel{
		db: goqu,
	}, nil
}

func (r *RatingModel) ListRatings(page, limit uint) ([]MovieRating, error) {
	var ratings []MovieRating

	offset := (page - 1) * limit

	ds := r.db.From(RatingsTable).
		Select(
			goqu.T(RatingsTable).Col("movie_id"),
			goqu.T(MovieTable).Col("title"),
			goqu.Func("AVG", goqu.T(RatingsTable).Col("rating")).As("avg_rating"),
		)

	ds = ds.Join(goqu.T(MovieTable), goqu.On(goqu.T(MovieTable).Col("id").Eq(goqu.T(RatingsTable).Col("movie_id")))).
		GroupBy("ratings.movie_id", "movies.title").
		Limit(limit).
		Offset(offset)

	err := ds.ScanStructs(&ratings)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ratings: %w", err)
	}

	return ratings, nil
}

func (r *RatingModel) GetRating(id string) (MovieRating, error) {
	var rating MovieRating

	movieID, err := strconv.Atoi(id)
	if err != nil {
		return MovieRating{}, fmt.Errorf("invalid movie ID: %w", err)
	}

	found, err := r.db.From(RatingsTable).
		Select(
			goqu.T(RatingsTable).Col("movie_id"),
			goqu.T(MovieTable).Col("title"),
			goqu.Func("AVG", goqu.T(RatingsTable).Col("rating")).As("avg_rating"),
		).
		Join(goqu.T(MovieTable), goqu.On(goqu.T(MovieTable).Col("id").Eq(goqu.T(RatingsTable).Col("movie_id")))).
		Where(goqu.T(RatingsTable).Col("movie_id").Eq(movieID)).
		GroupBy(goqu.T(RatingsTable).Col("movie_id"), goqu.T(MovieTable).Col("title")).
		ScanStruct(&rating)

	if err != nil {
		return MovieRating{}, fmt.Errorf("failed to fetch rating: %w", err)
	}

	if !found {
		return MovieRating{}, sql.ErrNoRows
	}

	return rating, nil
}

func (r *RatingModel) DeleteRatings(movieId, userId int) error {
	res, err := r.db.From(RatingsTable).Delete().
		Where(goqu.Ex{
			"user_id":  userId,
			"movie_id": movieId,
		}).
		Executor().
		Exec()

	if err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RatingModel) UpdateRatings(userId, movieId int, newRating float32) error {
	ds := r.db.Update(RatingsTable).Set(goqu.Record{
		"rating":    newRating,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}).Where(goqu.Ex{
		"user_id":  userId,
		"movie_id": movieId,
	})

	res, err := ds.Executor().Exec()
	if err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *RatingModel) AddorUpdateRatings(rating *Ratings) error {
	var count int
	_, err := r.db.From(MovieTable).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"id": rating.MovieId}).
		ScanVal(&count)

	if err != nil {
		return fmt.Errorf("failed to check movie existence: %w", err)
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	insert := r.db.Insert(RatingsTable).Rows(
		goqu.Record{
			"user_id":  rating.UserId,
			"movie_id": rating.MovieId,
			"rating":   rating.Rating,
		},
	).OnConflict(
		goqu.DoUpdate(
			"user_id, movie_id", goqu.Record{
				"rating":    rating.Rating,
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			}),
	)

	_, err = insert.Executor().Exec()
	if err != nil {
		return fmt.Errorf("failed to upsert rating: %w", err)
	}

	return nil
}
