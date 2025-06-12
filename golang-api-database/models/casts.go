package models

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/doug-martin/goqu/v9"
)

var CastTable = "movie_casts"

type MovieCast struct {
	MovieID   int    `db:"movie_id"`
	PersonID  int    `db:"person_id"`
	CreditID  string `db:"credit_id"`
	CastID    int    `db:"cast_id"`
	Character string `db:"character"`
	Name      string `db:"name"`
	Order     int    `db:"cast_order"`
}

type CastsModel struct {
	db *goqu.Database
}

func InitCastsModel(goqu *goqu.Database) (*CastsModel, error) {
	return &CastsModel{
		db: goqu,
	}, nil
}

func (c *CastsModel) ListCasts(id string) ([]MovieCast, error) {
	var casts []MovieCast

	movieID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	err = c.db.From(CastTable).
		Select("person_id", "movie_id", "name", "credit_id", "cast_id", "character", "cast_order").
		Join(goqu.T("credits"), goqu.On(goqu.T(CastTable).Col("person_id").Eq(goqu.T("credits").Col("id")))).
		Where(goqu.T(CastTable).Col("movie_id").Eq(movieID)).
		ScanStructs(&casts)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch cast: %w", err)
	}

	if len(casts) == 0 {
		return nil, sql.ErrNoRows
	}

	return casts, nil
}

type ActorWithMovies struct {
	ActorName string   `json:"actor"`
	Movies    []string `json:"movies"`
}

func (c *CastsModel) ListMoviesByCastId(id string) (*ActorWithMovies, error) {
	var actorName string
	var movies []string

	castID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid cast ID: %w", err)
	}

	dsActor := c.db.From("credits").Select("name").Where(goqu.C("id").Eq(castID))
	if found, err := dsActor.ScanVal(&actorName); err != nil || !found {
		return nil, sql.ErrNoRows
	}

	dsMovies := c.db.From(CastTable).Select("movies.title").
		Join(goqu.T(MovieTable), goqu.On(goqu.T(CastTable).Col("movie_id").Eq(goqu.T(MovieTable).Col("id")))).
		Where(goqu.T(CastTable).Col("person_id").Eq(castID))

	if err := dsMovies.ScanVals(&movies); err != nil {
		return nil, fmt.Errorf("failed to fetch movies for actor: %w", err)
	}

	return &ActorWithMovies{ActorName: actorName, Movies: movies}, nil
}

var ErrCastAlreadyExists = errors.New("cast already exists")

func (c *CastsModel) AddMovieCasts(cast *MovieCast) error {
	var moviecount int
	_, err := c.db.From(MovieTable).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"id": cast.MovieID}).
		ScanVal(&moviecount)

	if err != nil {
		return fmt.Errorf("failed to check movie existence: %w", err)
	}

	if moviecount == 0 {
		return ErrMovieNotFound
	}

	var creditcount int
	_, err = c.db.From(CreditsTable).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"id": cast.PersonID}).
		ScanVal(&creditcount)

	if err != nil {
		return fmt.Errorf("failed to check credit existence: %w", err)
	}

	if creditcount == 0 {
		return ErrCreditNotFound
	}

	insert := c.db.Insert(CastTable).Rows(
		goqu.Record{
			"movie_id":   cast.MovieID,
			"person_id":  cast.PersonID,
			"credit_id":  GenerateID(),
			"cast_id":    generateIntegerID(),
			"character":  cast.Character,
			"cast_order": cast.Order,
		},
	).OnConflict(goqu.DoNothing())

	var insertedId string
	_, err = insert.Returning("credit_id").Executor().ScanVal(&insertedId)
	if err != nil {
		return fmt.Errorf("failed to insert movie cast: %w", err)
	}
	if insertedId == "" {
		return ErrCastAlreadyExists
	}

	return nil
}

func generateIntegerID() int {
	return rand.Intn(100) + 1
}
