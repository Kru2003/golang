package controllers

import (
	"encoding/json"
	"net/http"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/constants"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/models"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CastController struct {
	castModel *models.CastModel
	logger    *zap.Logger
}

func NewCastController(logger *zap.Logger) (*CastController, error) {
	model := models.NewCastModel()
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
//	200: ResponseListCrewMembers
//	500: GenericErrorResponse
func (ctrl *CastController) ListCastMembers(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)
	castMembers, err := ctrl.castModel.ListCastMembers(movieId)
	if err != nil {
		ctrl.logger.Error(constants.LoadCreditsError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.LoadCreditsError)
	}

	return utils.JSONSuccess(c, http.StatusOK, castMembers)
}

// ListMoviesByCastId retrieves all movies IDs by castID
// swagger:route GET /actor/{castId}/cast Cast ListMoviesByCastId
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
//	500: GenericErrorResponse
func (ctrl *CastController) ListMoviesByCastId(c *fiber.Ctx) error {
	castId := c.Params(constants.CastId)
	movies, err := ctrl.castModel.ListMoviesByCastId(castId)
	if err != nil {
		ctrl.logger.Error(constants.LoadCreditsError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.LoadCreditsError)
	}

	return utils.JSONSuccess(c, http.StatusOK, movies)

}

// UpdateCastMember updates a cast member details
// swagger:route PUT /movies/{movieId}/cast/{castId} Cast UpdateCastMember
//
// Updates an existing cast member.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestUpdateCastMember
//
// Responses:
//
//	200: GenericSuccessResponse
//	400: ValidationErrorResponse
//	500: GenericErrorResponse
func (ctrl *CastController) UpdateCastMember(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)
	castId := c.Params(constants.CastId)

	var validate = validator.New()
	var cast models.CastMember

	// Unmarshal JSON data into the movie struct
	if err := json.Unmarshal(c.Body(), &cast); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	// Validate the movie struct using validator
	if err := validate.Struct(cast); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.castModel.UpdateCastMember(movieId, castId, cast); err != nil {
		ctrl.logger.Error(constants.UpdateCastError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UpdateCastError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.UpdateCastSuccess)
}
