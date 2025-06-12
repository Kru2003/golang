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

type CastController struct {
	castModel *models.CastsModel
	logger    *zap.Logger
}

func NewCastController(goqu *goqu.Database, logger *zap.Logger) (*CastController, error) {
	model, err := models.InitCastsModel(goqu)
	if err != nil {
		return nil, err
	}
	return &CastController{
		castModel: model,
		logger:    logger,
	}, nil
}

// ListCastMembers retrieves all cast members by movieID
// swagger:route GET /movies/{movieId}/casts Cast ListCastMembers
//
// Retrieves cast members of a movie by ID.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestListCastMembers
//
// Responses:
//
//	200: ResponseListCastMembers
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *CastController) ListCastMembers(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)

	casts, err := ctrl.castModel.ListCasts(movieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.CastsNotExist)
		}
		ctrl.logger.Error("error while get casts of movie by id", zap.Any("id", movieId), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetCasts)
	}

	return utils.JSONSuccess(c, http.StatusOK, casts)
}

// ListMoviesByCastId retrieves all movie titles in which actor had played role
// swagger:route GET /actor/{castId}/movies Cast ListMoviesByCastId
//
// Retrieves all movies by cast ID.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestListMoviesByCastId
//
// Responses:
//
//	200: ResponseListMoviesByCastId
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *CastController) ListMoviesByCastId(c *fiber.Ctx) error {
	castId := c.Params(constants.CastId)

	movies, err := ctrl.castModel.ListMoviesByCastId(castId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.ActorNotExist)
		}
		ctrl.logger.Error("error while get movies by cast id", zap.Any("id", castId), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetMovie)
	}

	return utils.JSONSuccess(c, http.StatusOK, movies)
}

// AddMovieCastsMember add a cast member of a movie
// swagger:route POST /movies/{movieId}/credit/{creditId}/cast Cast AddMovieCastMember
//
// Adds a cast member to a movie.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestAddMovieCastMember
//
// Responses:
//
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *CastController) AddMovieCastMember(c *fiber.Ctx) error {
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
		Character string `json:"character" validate:"required"`
		Order     int    `json:"order" validate:"required"`
	}

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	cast := models.MovieCast{
		MovieID:   movieid,
		PersonID:  creditid,
		Character: input.Character,
		Order:     input.Order,
	}

	validate := validator.New()
	if err := validate.Struct(cast); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.castModel.AddMovieCasts(&cast); err != nil {
		if errors.Is(err, models.ErrMovieNotFound) || errors.Is(err, models.ErrCreditNotFound) {
			return utils.JSONFail(c, http.StatusNotFound, "Movie or credit not found")
		}

		if errors.Is(err, models.ErrCastAlreadyExists) {
			return utils.JSONFail(c, http.StatusBadRequest, "movie cast already exists")
		}

		ctrl.logger.Error(constants.ErrAddMovieCast, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrAddMovieCast)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AddMovieCastSuccess)
}
