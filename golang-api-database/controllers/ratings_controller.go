package controllers

import (
	"database/sql"
	"encoding/json"
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

// RatingsController for ratingModel and movieModel controllers
type RatingsController struct {
	ratingModel *models.RatingModel
	logger      *zap.Logger
}

// NewRatingsController is to initialize RatingsController
func NewRatingsController(goqu *goqu.Database, logger *zap.Logger) (*RatingsController, error) {
	model, err := models.InitRatingsModel(goqu)
	if err != nil {
		return nil, err
	}
	return &RatingsController{
		ratingModel: model,
		logger:      logger,
	}, nil
}

// ListAllMovieRatings displays average ratings of all movies
// swagger:route GET /ratings/movies Ratings ListAllMovieRatings
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
//	401: GenericResFailBadRequest
//	500: GenericResError
func (ctrl *RatingsController) ListAllMovieRatings(c *fiber.Ctx) error {
	page, limit, err := PaginationQuery(c)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidPageOrLimit)
	}

	ratings, err := ctrl.ratingModel.ListRatings(page, limit)
	if err != nil {
		ctrl.logger.Error(constants.ErrGetRatings, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetRatings)
	}

	return utils.JSONSuccess(c, http.StatusOK, ratings)
}

// GetRatingsByMovieId retrieves ratings of a movie by its ID
// swagger:route GET /movies/{movieId}/ratings Ratings GetRatingsByMovieId
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
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *RatingsController) GetRatingByMovieId(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)

	rating, err := ctrl.ratingModel.GetRating(movieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.MovieNotExist)
		}
		ctrl.logger.Error("error while get rating of movie by id", zap.Any("id", movieId), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetRatings)
	}

	return utils.JSONSuccess(c, http.StatusOK, rating)
}

// DeleteRating deletes ratings of a movie
// swagger:route DELETE /ratings/movies/{movieId}/user/{userId}/ratings Ratings DeleteRating
//
// Deletes rating of a movie from the system.
// //
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
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *RatingsController) DeleteRating(c *fiber.Ctx) error {
	userId := c.Params(constants.UserId)
	movieId := c.Params(constants.ParamMid)

	movieid, err := strconv.Atoi(movieId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "movie ID must be a valid integer")
	}

	userid, err := strconv.Atoi(userId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "user ID must be a valid integer")
	}
	err = ctrl.ratingModel.DeleteRatings(movieid, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.RatingNotExist)
		}
		ctrl.logger.Error(constants.ErrDeleteRating, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrDeleteRating)
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
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *RatingsController) UpdateRating(c *fiber.Ctx) error {
	userId := c.Params(constants.UserId)
	movieId := c.Params(constants.ParamMid)

	movieid, err := strconv.Atoi(movieId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "movie ID must be a valid integer")
	}

	userid, err := strconv.Atoi(userId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "user ID must be a valid integer")
	}

	var updateData struct {
		Rating float32 `json:"rating"`
	}

	if err := json.Unmarshal(c.Body(), &updateData); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	err = ctrl.ratingModel.UpdateRatings(userid, movieid, updateData.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.MovieNotExist)
		}
		ctrl.logger.Error(constants.ErrUpdateRating, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUpdateRating)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.UpdateRatingSuccess)
}

// AddRating adds ratings of a movie
// swagger:route POST /ratings/user/{userId}/ratings Ratings AddRating
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
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	500: GenericResError
func (ctrl *RatingsController) AddRating(c *fiber.Ctx) error {
	userId := c.Params(constants.UserId)

	userid, err := strconv.Atoi(userId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "user ID must be a valid integer")
	}

	// Parse only movieId and rating from body
	var input struct {
		MovieId int     `json:"movieId" validate:"required"`
		Rating  float32 `json:"rating" validate:"required"`
	}

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	// Construct the Ratings struct
	rating := models.Ratings{
		UserId:  userid,
		MovieId: input.MovieId,
		Rating:  input.Rating,
	}

	validate := validator.New()
	if err := validate.Struct(rating); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ValidationFailed)
	}

	if err := ctrl.ratingModel.AddorUpdateRatings(&rating); err != nil {
		ctrl.logger.Error(constants.ErrAddRating, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrAddRating)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AddRatingSuccess)
}
