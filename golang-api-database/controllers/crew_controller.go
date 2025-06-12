package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/constants"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/models"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CrewController struct {
	crewModel *models.CrewModel
	logger    *zap.Logger
}

func NewCrewController(goqu *goqu.Database, logger *zap.Logger) (*CrewController, error) {
	model, err := models.InitCrewModel(goqu)
	if err != nil {
		return nil, err
	}
	return &CrewController{
		crewModel: model,
		logger:    logger,
	}, nil
}

// ListCrewMembers retrieves all cast members by movieID
// swagger:route GET /movies/{movieId}/crew Crew ListCrewMembers
//
// Retrieves crew members of a movie by ID.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestListCrewMembers
//
// Responses:
//
//	200: ResponseListCrewMembers
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *CrewController) ListCrewMembers(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)

	crew, err := ctrl.crewModel.ListCrew(movieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.MovieNotExist)
		}
		ctrl.logger.Error("error while get crew of movie by id", zap.Any("id", movieId), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetCrew)
	}

	return utils.JSONSuccess(c, http.StatusOK, crew)
}

// AddMovieCrewMember add a crew member of a movie
// swagger:route POST /movies/{movieId}/credit/{creditId}/crew Crew AddMovieCrewMember
//
// Adds a crew member to a movie.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestAddMovieCrewMember
//
// Responses:
//
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *CrewController) AddMovieCrewMember(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)
	creditId := c.Params(constants.CreditId)

	movieid, err := strconv.Atoi(movieId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "movie ID must be a valid integer")
	}

	creditid, err := strconv.Atoi(creditId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "credit ID must be a valid integer")
	}

	var input struct {
		Department string `json:"department" validate:"required"`
		Job        string `json:"job" validate:"required"`
	}

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	crew := models.MovieCrew{
		MovieID:    movieid,
		PersonID:   creditid,
		Department: input.Department,
		Job:        input.Job,
	}

	validate := validator.New()
	if err := validate.Struct(crew); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.crewModel.AddMovieCrew(&crew); err != nil {
		if errors.Is(err, models.ErrMovieNotFound) || errors.Is(err, models.ErrCreditNotFound) {
			return utils.JSONFail(c, http.StatusNotFound, "Movie or credit not found")
		}

		if errors.Is(err, models.ErrCrewAlreadyExists) {
			return utils.JSONFail(c, http.StatusBadRequest, "movie crew already exists")
		}

		ctrl.logger.Error(constants.ErrAddMovieCrew, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrAddMovieCrew)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AddMovieCrewSuccess)
}
