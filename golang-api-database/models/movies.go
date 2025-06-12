package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/go-playground/validator"
)

// MovieTable represent table name
const MovieTable = "movies"

type MovieDB struct {
	ID               int             `db:"id"`
	IMDB_ID          sql.NullString  `db:"imdb_id"`
	OriginalTitle    string          `db:"original_title"`
	OriginalLanguage string          `db:"original_language"`
	Title            string          `db:"title"`
	Tagline          sql.NullString  `db:"tagline"`
	Overview         sql.NullString  `db:"overview"`
	Popularity       float64         `db:"popularity"`
	Status           sql.NullString  `db:"status"`
	ReleaseDate      sql.NullString  `db:"release_date"`
	Runtime          sql.NullFloat64 `db:"runtime"`
	Vote_average     float64         `db:"vote_average"`
	Vote_count       int64           `db:"vote_count"`
}

type Movie struct {
	ID               int     `json:"id"`
	IMDB_ID          string  `json:"imdb_id"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Title            string  `json:"title"`
	Tagline          string  `json:"tagline,omitempty"`
	Overview         string  `json:"overview,omitempty"`
	Popularity       float64 `json:"popularity"`
	Status           string  `json:"status"`
	ReleaseDate      string  `json:"release_date,omitempty"`
	Runtime          float64 `json:"runtime"`
	Vote_average     float64 `json:"vote_average"`
	Vote_count       int64   `json:"vote_count"`
}

func ConvertMovieDBToMovie(m MovieDB) Movie {
	return Movie{
		ID:               m.ID,
		IMDB_ID:          nullStringToString(m.IMDB_ID),
		OriginalTitle:    m.OriginalTitle,
		OriginalLanguage: m.OriginalLanguage,
		Title:            m.Title,
		Tagline:          nullStringToString(m.Tagline),
		Overview:         nullStringToString(m.Overview),
		Popularity:       m.Popularity,
		Status:           nullStringToString(m.Status),
		ReleaseDate:      nullStringToString(m.ReleaseDate),
		Runtime:          nullFloatToFloat(m.Runtime),
		Vote_average:     m.Vote_average,
		Vote_count:       m.Vote_count,
	}
}

func nullFloatToFloat(ns sql.NullFloat64) float64 {
	if ns.Valid {
		return ns.Float64
	}
	return 0
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

type MovieModel struct {
	db *goqu.Database
}

func InitMovieModel(goqu *goqu.Database) (*MovieModel, error) {
	return &MovieModel{
		db: goqu,
	}, nil
}

func (m *MovieModel) GetMovie(id string) (Movie, error) {
	var movieDB MovieDB
	found, err := m.db.From(MovieTable).Where(goqu.Ex{"id": id}).
		Select("id", "imdb_id", "original_title", "original_language", "title", "status", "vote_average", "vote_count", "popularity", "release_date", "tagline", "overview", "runtime").
		ScanStruct(&movieDB)
	if err != nil {
		return Movie{}, err
	}

	if !found {
		return Movie{}, sql.ErrNoRows
	}

	return ConvertMovieDBToMovie(movieDB), nil
}

func (m *MovieModel) ListMovies(filters map[string]string, page, limit uint) ([]Movie, error) {
	var movieDBs []MovieDB
	var movies []Movie

	ds := m.db.From(MovieTable).
		Select(goqu.DISTINCT("movies.id"), "imdb_id", "original_language", "original_title", "title", "status", "vote_average", "vote_count", "popularity", "release_date", "tagline", "overview", "runtime")

	ds = ds.Join(goqu.T("movie_genres"), goqu.On(goqu.T("movies").Col("id").Eq(goqu.T("movie_genres").Col("movieid")))).
		Join(goqu.T("genres"), goqu.On(goqu.T("movie_genres").Col("genreid").Eq(goqu.T("genres").Col("id")))).
		Join(goqu.T("movie_languages"), goqu.On(goqu.T(MovieTable).Col("id").Eq(goqu.T("movie_languages").Col("movieid"))))

	if languageFilter, ok := filters["language"]; ok && languageFilter != "" {
		ds = ds.Where(goqu.T("movie_languages").Col("language_code").Eq(languageFilter))
	}

	if genreFilter, ok := filters["genre"]; ok && genreFilter != "" {
		ds = ds.Where(goqu.L("LOWER(?) = LOWER(?)", goqu.T("genres").Col("name"), genreFilter))
	}

	if nameFilter, ok := filters["name"]; ok && nameFilter != "" {
		ds = ds.Where(goqu.T(MovieTable).Col("original_title").ILike("%" + nameFilter + "%"))
	}

	ds = ds.Offset((page - 1) * limit).Limit(limit)

	// Scan into MovieDB structs
	err := ds.ScanStructs(&movieDBs)
	if err != nil {
		return nil, fmt.Errorf("error fetching movies: %w", err)
	}

	// Convert each MovieDB to Movie
	for _, mdb := range movieDBs {
		movies = append(movies, ConvertMovieDBToMovie(mdb))
	}

	return movies, nil
}

func (m *MovieModel) DeleteMovie(id int) error {
	res, err := m.db.From(MovieTable).
		Delete().
		Where(goqu.Ex{"id": id}).
		Executor().
		Exec()

	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func ValidateReleaseDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

type MovieWithMetadata struct {
	OriginalTitle    string   `json:"original_title" validate:"required"`
	OriginalLanguage string   `json:"original_language" validate:"required"`
	Title            string   `json:"title" validate:"required"`
	Overview         string   `json:"overview"`
	Popularity       float64  `json:"popularity" validate:"required,gte=0"`
	Status           string   `json:"status"`
	ReleaseDate      string   `json:"release_date" validate:"required,releaseDateFormat"`
	Runtime          float64  `json:"runtime" validate:"gte=0"`
	Vote_average     float64  `json:"vote_average" validate:"required,gte=0"`
	Vote_count       int64    `json:"vote_count" validate:"required,gte=0"`
	Genres           []string `json:"genres"`
	Languages        []string `json:"languages"`
}

func getNextID(tx *goqu.TxDatabase, tableName string) (int64, error) {
	var maxID int64
	found, err := tx.From(tableName).Select(goqu.MAX("id")).ScanVal(&maxID)
	if err != nil {
		return 0, err
	}
	if !found {
		return 1, nil
	}
	return maxID + 1, nil
}

func (m *MovieModel) AddMovie(movie *MovieWithMetadata) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	movieID, err := getNextID(tx, MovieTable)
	if err != nil {
		return fmt.Errorf("failed to get next movie ID: %w", err)
	}

	_, err = tx.Insert(MovieTable).Rows(goqu.Record{
		"id":                movieID,
		"original_title":    movie.OriginalTitle,
		"original_language": movie.OriginalLanguage,
		"title":             movie.Title,
		"overview":          movie.Overview,
		"popularity":        movie.Popularity,
		"status":            movie.Status,
		"release_date":      movie.ReleaseDate,
		"runtime":           movie.Runtime,
		"vote_average":      movie.Vote_average,
		"vote_count":        movie.Vote_count,
	}).Executor().Exec()

	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	err = m.handleGenres(tx, movieID, movie.Genres)
	if err != nil {
		return err
	}

	err = m.handleLanguages(tx, movieID, movie.Languages)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (m *MovieModel) handleGenres(tx *goqu.TxDatabase, movieID int64, genres []string) error {
	for _, g := range genres {
		var genreID int64
		found, err := tx.From("genres").
			Where(goqu.L("LOWER(name) = LOWER(?)", g)).
			Select("id").ScanVal(&genreID)
		if err != nil {
			return fmt.Errorf("error querying genre: %w", err)
		}
		if !found {
			genreID, err = getNextID(tx, "genres")
			if err != nil {
				return fmt.Errorf("failed to get next genre ID: %w", err)
			}
			_, err = tx.Insert("genres").Rows(goqu.Record{"id": genreID, "name": g}).Executor().Exec()
			if err != nil {
				return fmt.Errorf("failed to insert genre: %w", err)
			}
		}
		_, err = tx.Insert("movie_genres").Rows(goqu.Record{
			"movieid": movieID,
			"genreid": genreID,
		}).Executor().Exec()
		if err != nil {
			return fmt.Errorf("failed to link movie and genre: %w", err)
		}
	}
	return nil
}

func (m *MovieModel) handleLanguages(tx *goqu.TxDatabase, movieID int64, languages []string) error {
	for _, isoCode := range languages {
		var lang string
		found, err := tx.From("languages").
			Where(goqu.C("iso_code").Eq(isoCode)).
			Select("iso_code").ScanVal(&lang)
		if err != nil {
			return fmt.Errorf("error checking language: %w", err)
		}

		if !found {
			_, err := tx.Insert("languages").Rows(goqu.Record{
				"iso_code": isoCode,
			}).Executor().Exec()
			if err != nil {
				return fmt.Errorf("failed to insert language: %w", err)
			}
		}

		_, err = tx.Insert("movie_languages").Rows(goqu.Record{
			"movieid":       movieID,
			"language_code": isoCode,
		}).Executor().Exec()
		if err != nil {
			return fmt.Errorf("failed to link movie and language: %w", err)
		}
	}
	return nil
}

func (m *MovieModel) UpdateMovie(movieID int, movie *MovieWithMetadata) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row, err := tx.Update(MovieTable).
		Set(goqu.Record{
			"original_title":    movie.OriginalTitle,
			"original_language": movie.OriginalLanguage,
			"title":             movie.Title,
			"overview":          movie.Overview,
			"popularity":        movie.Popularity,
			"status":            movie.Status,
			"release_date":      movie.ReleaseDate,
			"runtime":           movie.Runtime,
			"vote_average":      movie.Vote_average,
			"vote_count":        movie.Vote_count,
		}).
		Where(goqu.C("id").Eq(movieID)).
		Executor().Exec()
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}
	if rowsAffected, _ := row.RowsAffected(); rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if _, err = tx.Delete("movie_genres").Where(goqu.C("movieid").Eq(movieID)).Executor().Exec(); err != nil {
		return fmt.Errorf("failed to delete old genres: %w", err)
	}
	if err = m.handleGenres(tx, int64(movieID), movie.Genres); err != nil {
		return err
	}

	if _, err = tx.Delete("movie_languages").Where(goqu.C("movieid").Eq(movieID)).Executor().Exec(); err != nil {
		return fmt.Errorf("failed to delete old languages: %w", err)
	}
	if err = m.handleLanguages(tx, int64(movieID), movie.Languages); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
