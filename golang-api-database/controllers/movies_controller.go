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

// MovieController for movieModel controllers
type MovieController struct {
	movieModel *models.MovieModel
	logger     *zap.Logger
}

// NewMovieController is to intialize MovieController
func NewMovieController(goqu *goqu.Database, logger *zap.Logger) (*MovieController, error) {
	model, err := models.InitMovieModel(goqu)
	if err != nil {
		return nil, err
	}
	return &MovieController{
		movieModel: model,
		logger:     logger,
	}, nil
}

// PaginationQuery is to handle page and limit query
func PaginationQuery(c *fiber.Ctx) (uint, uint, error) {
	pageStr := c.Query("page", "1")    // Default: 1
	limitStr := c.Query("limit", "10") // Default: 10

	// Convert query params to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, err
	}
	return uint(page), uint(limit), nil
}

// GetMovieByID retrieves a movie by ID
// swagger:route GET /movies/{movieId} Movies GetMovieByID
//
// Retrieves movie details by ID.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestGetMovieByID
//
// Responses:
//
//	200: ResponseGetMovieByID
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *MovieController) GetMovieByID(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)
	movie, err := ctrl.movieModel.GetMovie(movieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.MovieNotExist)
		}
		ctrl.logger.Error("error while get user by id", zap.Any("id", movieId), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetMovie)
	}

	return utils.JSONSuccess(c, http.StatusOK, movie)
}

// ListMovies lists all movies with pagination
// swagger:route GET /movies Movies ListMovies
//
// Retrieves a paginated list of movies.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestListMovies
//
// Responses:
//
//	200: ResponseListMovies
//	401: GenericResFailBadRequest
//	500: GenericResError
func (ctrl *MovieController) ListMovies(c *fiber.Ctx) error {
	// Extract query parameters
	filters := map[string]string{
		"name":     c.Query("name"),
		"genre":    c.Query("genre"),
		"language": c.Query("language"),
	}

	page, limit, err := PaginationQuery(c)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidPageOrLimit)
	}

	movies, err := ctrl.movieModel.ListMovies(filters, page, limit)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetMovie)
	}

	return utils.JSONSuccess(c, http.StatusOK, movies)
}

// DeleteMovieById deletes a movie by ID
// swagger:route DELETE /movies/{movieId} Movies DeleteMovieById
//
// Deletes a movie from the system.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestDeleteMovieByID
//
// Responses:
//
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *MovieController) DeleteMovieById(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)

	id, err := strconv.Atoi(movieId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "movie ID must be a valid integer")
	}

	err = ctrl.movieModel.DeleteMovie(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.MovieNotExist)
		}
		ctrl.logger.Error(constants.ErrDeleteMovie, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrDeleteMovie)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.DeleteMovieSuccess)
}

// AddMovie adds a new movie
// swagger:route POST /movies Movies AddMovie
//
// Adds a new movie to the system.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestAddMovie
//
// Responses:
//
//	200: GenericResOk
//	400: GenericResFailBadRequest
//	500: GenericResError
func (ctrl *MovieController) AddMovie(c *fiber.Ctx) error {
	var validate = validator.New()
	var movie models.MovieWithMetadata

	// Unmarshal JSON data into the movie struct
	if err := json.Unmarshal(c.Body(), &movie); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	validate.RegisterValidation("releaseDateFormat", models.ValidateReleaseDate)
	// Validate the movie struct using validator
	if err := validate.Struct(movie); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.movieModel.AddMovie(&movie); err != nil {
		ctrl.logger.Error(constants.ErrAddMovie, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrAddMovie)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AddMovieSuccess)
}

// UpdateMovie updates a movie by ID
// swagger:route PUT /movies/{movieId} Movies UpdateMovie
//
// Updates an existing movie.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestUpdateMovie
//
// Responses:
//
//	200: GenericResOk
//	401: GenericResFailBadRequest
//	404: GenericResFailNotFound
//	500: GenericResError
func (ctrl *MovieController) UpdateMovie(c *fiber.Ctx) error {
	movieId := c.Params(constants.ParamMid)
	id, err := strconv.Atoi(movieId)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "movie ID must be a valid integer")
	}

	var validate = validator.New()
	var movie models.MovieWithMetadata

	// Unmarshal JSON data into the movie struct
	if err := json.Unmarshal(c.Body(), &movie); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	validate.RegisterValidation("releaseDateFormat", models.ValidateReleaseDate)
	// Validate the movie struct using validator
	if err := validate.Struct(movie); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.movieModel.UpdateMovie(id, &movie); err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.MovieNotExist)
		}
		ctrl.logger.Error(constants.UpdateMovieError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UpdateMovieError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.UpdateMovieSuccess)
}
