package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	constants "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/constants"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/models"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RatingsController for ratingModel and movieModel controllers
type RatingsController struct {
	ratingModel *models.RatingModel
	movieModel  *models.MovieModel
	logger      *zap.Logger
}

// NewRatingsController is to initialize RatingsController
func NewRatingsController(logger *zap.Logger) (*RatingsController, error) {
	model := models.NewRatingsModel()
	movieModel := models.NewMovieModel()
	return &RatingsController{
		ratingModel: model,
		movieModel:  movieModel,
		logger:      logger,
	}, nil
}

// ListAllMovieRatings displays average ratings of all movies
// swagger:route GET /ratings Ratings ListAllMovieRatings
//
// Retrieves a paginated ratings of movies.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestListAllMovieRatings
//
// Responses:
//
//	200: ResponseListAllMovieRatings
//	500: GenericErrorResponse
func (ctrl *RatingsController) ListAllMovieRatings(c *fiber.Ctx) error {
	page, limit, err := PaginationQuery(c)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.InvalidPageOrLimitError)
	}

	ratings, err := ctrl.ratingModel.ListRatings(page, limit)
	if err != nil {
		ctrl.logger.Error(constants.LoadRatingsError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.LoadRatingsError)
	}

	return utils.JSONSuccess(c, http.StatusOK, ratings)
}

// GetRatingsByMovieId retrieves ratings of a movie by its ID
// swagger:route GET /ratings/movies/{movieId}/ratings Ratings GetRatingsByMovieId
//
// Retrieves a ratings of movie by ID.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestGetRatingsByMovieId
//
// Responses:
//
//	200: ResponseGetRatingsByMovieId
//	500: GenericErrorResponse
func (ctrl *RatingsController) GetRatingsByMovieId(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)

	ratings, err := ctrl.ratingModel.GetRatingsByMovieId(movieId)
	if err != nil {
		ctrl.logger.Error(constants.LoadRatingsError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.LoadRatingsError)
	}

	return utils.JSONSuccess(c, http.StatusOK, ratings)
}

// AddRating adds ratings of a movie
// swagger:route POST /ratings Ratings AddRating
//
// Adds a new ratings of a movie to the system.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestAddRating
//

// Responses:
//
//	200: GenericSuccessResponse
//	400: ValidationErrorResponse
//	500: GenericErrorResponse
func (ctrl *RatingsController) AddRating(c *fiber.Ctx) error {
	var rating models.Ratings
	var validate = validator.New()

	// Read and parse request body
	if err := json.Unmarshal(c.Body(), &rating); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	// Validate rating fields
	if err := validate.Struct(rating); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	// Check if the movie exists in the database or file
	exists, err := ctrl.movieModel.MovieExists(rating.MovieId)
	if err != nil {
		ctrl.logger.Error(constants.MovieCheckError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.MovieCheckError)
	}

	if !exists {
		return utils.JSONError(c, http.StatusNotFound, constants.MovieCheckError)
	}

	// Save rating if movie exists
	if err := ctrl.ratingModel.AddRatings(&rating); err != nil {
		ctrl.logger.Error(constants.AddRatingError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.AddRatingError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AddRatingSuccess)
}

// DeleteRating deletes ratings of a movie
// swagger:route DELETE /ratings/movies/{movieId}/user/{userId}/ratings Ratings DeleteRating
//
// Deletes rating of a movie from the system.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestDeleteRating
//
// Responses:
//
//	200: GenericSuccessResponse
//	400: GenericErrorResponse
//	500: GenericErrorResponse
func (ctrl *RatingsController) DeleteRating(c *fiber.Ctx) error {
	userId := c.Params(constants.UserId)
	movieId := c.Params(constants.MovieId)

	if movieId == "" && userId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.ValidationFailed)
	}

	err := ctrl.ratingModel.DeleteRatings(movieId, &userId)
	if err != nil {
		ctrl.logger.Error(constants.DeleteRatingError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.DeleteRatingError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.DeleteRatingSuccess)
}

// UpdateRating updates rating and its timestamp
// swagger:route PUT /ratings/movies/{movieId}/user/{userId}/ratings Ratings UpdateRating
//
// Updates an existing rating.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestUpdateRating
//
// Responses:
//
//	200: GenericSuccessResponse
//	400: ValidationErrorResponse
//	500: GenericErrorResponse
func (ctrl *RatingsController) UpdateRating(c *fiber.Ctx) error {
	userId := c.Params(constants.UserId)
	movieId := c.Params(constants.MovieId)

	if movieId == "" && userId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	var updateData struct {
		Rating string `json:"rating"`
	}

	if err := json.Unmarshal(c.Body(), &updateData); err != nil {
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	newTimestamp := fmt.Sprintf("%d", time.Now().Unix())

	err := ctrl.ratingModel.UpdateRatings(userId, movieId, updateData.Rating, newTimestamp)
	if err != nil {
		ctrl.logger.Error(constants.UpdateRatingError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UpdateRatingError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.UpdateRatingSuccess)

}
