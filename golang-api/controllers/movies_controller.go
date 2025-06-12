package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	constants "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/constants"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/models"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// MovieController for movieModel controllers
type MovieController struct {
	movieModel *models.MovieModel
	logger     *zap.Logger
}

// NewMovieController is to intialize MovieController
func NewMovieController(logger *zap.Logger) (*MovieController, error) {
	model := models.NewMovieModel()
	return &MovieController{
		movieModel: model,
		logger:     logger,
	}, nil
}

// PaginationQuery is to handle page and limit query
func PaginationQuery(c *fiber.Ctx) (int, int, error) {
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
	return page, limit, nil
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
//	500: GenericErrorResponse
func (ctrl *MovieController) ListMovies(c *fiber.Ctx) error {
	// Extract query parameters
	filters := map[string]string{
		"name":     c.Query("name"),
		"genre":    c.Query("genre"),
		"language": c.Query("language"),
	}

	page, limit, err := PaginationQuery(c)
	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, constants.InvalidPageOrLimitError)
	}

	movies, err := ctrl.movieModel.ListMovies(filters, page, limit)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, constants.PaginationError)
	}

	return utils.JSONSuccess(c, http.StatusOK, movies)
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
//	500: GenericErrorResponse
func (ctrl *MovieController) GetMovieByID(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)
	movie, err := ctrl.movieModel.GetMovie(movieId)
	if err != nil {
		ctrl.logger.Error(constants.LoadMoviesError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.LoadMoviesError)
	}

	return utils.JSONSuccess(c, http.StatusOK, movie)
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
//	200: GenericSuccessResponse
//	400: ValidationErrorResponse
//	500: GenericErrorResponse
func (ctrl *MovieController) AddMovie(c *fiber.Ctx) error {
	var validate = validator.New()
	var movie models.Movies

	// Unmarshal JSON data into the movie struct
	if err := json.Unmarshal(c.Body(), &movie); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	// Validate the movie struct using validator
	if err := validate.Struct(movie); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.movieModel.AddMovie(&movie); err != nil {
		ctrl.logger.Error(constants.AddMovieError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.AddMovieError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.AddMovieSuccess)
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
//	200: GenericSuccessResponse
//	400: GenericErrorResponse
//	500: GenericErrorResponse
func (ctrl *MovieController) DeleteMovieById(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)

	if movieId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	err := ctrl.movieModel.DeleteMovie(movieId)
	if err != nil {
		ctrl.logger.Error(constants.DeleteMovieError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.DeleteMovieError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.DeleteMovieSuccess)
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
//	200: GenericSuccessResponse
//	400: ValidationErrorResponse
//	500: GenericErrorResponse
func (ctrl *MovieController) UpdateMovie(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)

	var validate = validator.New()
	var movie models.Movies

	// Unmarshal JSON data into the movie struct
	if err := json.Unmarshal(c.Body(), &movie); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	if movie.ID != "" && movie.ID != movieId {
		ctrl.logger.Error("Movie ID mismatch")
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	// Validate the movie struct using validator
	if err := validate.Struct(movie); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.movieModel.UpdateMovie(movieId, &movie); err != nil {
		ctrl.logger.Error(constants.UpdateMovieError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UpdateMovieError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.UpdateMovieSuccess)
}
