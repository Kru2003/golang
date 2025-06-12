package routes

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/constants"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/controllers"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/middlewares"
	pMetrics "git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/pkg/prometheus"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, logger *zap.Logger, config config.AppConfig, pMetrics *pMetrics.PrometheusMetrics) error {
	mu.Lock()

	app.Use(middlewares.LogHandler(logger, pMetrics))

	app.Use(swagger.New(swagger.Config{
		FilePath: "./assets/swagger.json",
		Title:    "Swagger API Docs",
	}))

	err := setupMoviesController(app, logger)
	if err != nil {
		return err
	}

	err = setupRatingsController(app, logger)
	if err != nil {
		return err
	}

	err = setupCrewController(app, logger)
	if err != nil {
		return err
	}

	err = setupCastController(app, logger)
	if err != nil {
		return err
	}

	err = metricsController(app, logger, pMetrics)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func metricsController(app *fiber.App, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error {
	metricsController, err := controllers.InitMetricsController(logger, pMetrics)
	if err != nil {
		return nil
	}

	app.Get("/metrics", metricsController.Metrics)
	return nil
}

func setupMoviesController(app *fiber.App, logger *zap.Logger) error {
	movieController, err := controllers.NewMovieController(logger)
	if err != nil {
		logger.Error("Failed to initialize MovieController", zap.Error(err))
		return err
	}

	// Register movie routes
	movieRouter := app.Group("/movies")
	movieRouter.Get("/", movieController.ListMovies)
	movieRouter.Get(fmt.Sprintf("/:%s", constants.MovieId), movieController.GetMovieByID)
	movieRouter.Post("/", movieController.AddMovie)
	movieRouter.Delete(fmt.Sprintf("/:%s", constants.MovieId), movieController.DeleteMovieById)
	movieRouter.Put(fmt.Sprintf("/:%s", constants.MovieId), movieController.UpdateMovie)

	return nil

}

func setupRatingsController(app *fiber.App, logger *zap.Logger) error {
	ratingController, err := controllers.NewRatingsController(logger)
	if err != nil {
		logger.Error("Failed to intialize RatingController", zap.Error(err))
		return err
	}

	ratingRouter := app.Group("/ratings")

	ratingRouter.Get("/", ratingController.ListAllMovieRatings)
	ratingRouter.Get(fmt.Sprintf("/movies/:%s/ratings", constants.MovieId), ratingController.GetRatingsByMovieId)
	ratingRouter.Post("/", ratingController.AddRating)
	ratingRouter.Delete(fmt.Sprintf("/movies/:%s/user/:%s/ratings", constants.MovieId, constants.UserId), ratingController.DeleteRating)
	ratingRouter.Put(fmt.Sprintf("/movies/:%s/user/:%s/ratings", constants.MovieId, constants.UserId), ratingController.UpdateRating)

	return nil
}

func setupCastController(app *fiber.App, logger *zap.Logger) error {
	castController, err := controllers.NewCastController(logger)
	if err != nil {
		logger.Error("Failed to intialize CastController", zap.Error(err))
		return err
	}

	app.Get(fmt.Sprintf("/movies/:%s/casts", constants.MovieId), castController.ListCastMembers)
	app.Get(fmt.Sprintf("/actor/:%s/cast", constants.CastId), castController.ListMoviesByCastId)
	app.Put(fmt.Sprintf("/movies/:%s/casts/:%s", constants.MovieId, constants.CastId), castController.UpdateCastMember)

	return nil
}

func setupCrewController(app *fiber.App, logger *zap.Logger) error {
	crewController, err := controllers.NewCrewController(logger)
	if err != nil {
		logger.Error("Failed to intialize CrewController", zap.Error(err))
		return err
	}

	app.Get(fmt.Sprintf("/movies/:%s/crew", constants.MovieId), crewController.ListCrewMembers)
	app.Put(fmt.Sprintf("/movies/:%s/crew/:%s", constants.MovieId, constants.CrewId), crewController.UpdateCrewMember)

	return nil
}
