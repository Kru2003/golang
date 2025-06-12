package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CrewTable = "movie_crew"
var CreditsTable = "credits"

type MovieCrew struct {
	MovieID    int    `db:"movie_id" validate:"required"`
	PersonID   int    `db:"person_id" validate:"required"`
	CreditID   string `db:"credit_id"`
	Name       string `db:"name"`
	Department string `db:"department" validate:"required"`
	Job        string `db:"job" validate:"required"`
}

type CrewModel struct {
	db *goqu.Database
}

func InitCrewModel(goqu *goqu.Database) (*CrewModel, error) {
	return &CrewModel{
		db: goqu,
	}, nil
}

func (c *CrewModel) ListCrew(id string) ([]MovieCrew, error) {
	var crew []MovieCrew

	movieID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid movie ID: %w", err)
	}

	err = c.db.From(CrewTable).
		Select("person_id", "movie_id", "name", "credit_id", "job", "department").
		Join(goqu.T(CreditsTable), goqu.On(goqu.T(CrewTable).Col("person_id").Eq(goqu.T(CreditsTable).Col("id")))).
		Where(goqu.T(CrewTable).Col("movie_id").Eq(movieID)).
		ScanStructs(&crew)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch cast: %w", err)
	}

	if len(crew) == 0 {
		return nil, sql.ErrNoRows
	}

	return crew, nil
}

func GenerateID() string {
	return primitive.NewObjectID().Hex()
}

var (
	ErrMovieNotFound     = errors.New("movie not found")
	ErrCreditNotFound    = errors.New("credit not found")
	ErrCrewAlreadyExists = errors.New("crew entry already exists")
)

func (c *CrewModel) AddMovieCrew(crew *MovieCrew) error {
	var moviecount int
	_, err := c.db.From(MovieTable).
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"id": crew.MovieID}).
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
		Where(goqu.Ex{"id": crew.PersonID}).
		ScanVal(&creditcount)

	if err != nil {
		return fmt.Errorf("failed to check credit existence: %w", err)
	}

	if creditcount == 0 {
		return ErrCreditNotFound
	}

	insert := c.db.Insert(CrewTable).Rows(
		goqu.Record{
			"movie_id":   crew.MovieID,
			"person_id":  crew.PersonID,
			"credit_id":  GenerateID(),
			"department": crew.Department,
			"job":        crew.Job,
		},
	).OnConflict(goqu.DoNothing())

	var insertedId string
	_, err = insert.Returning("credit_id").Executor().ScanVal(&insertedId)
	if err != nil {
		return fmt.Errorf("failed to insert movie crew: %w", err)
	}
	if insertedId == "" {
		return ErrCrewAlreadyExists
	}

	return nil
}
