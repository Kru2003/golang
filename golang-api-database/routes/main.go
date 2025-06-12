package routes

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/config"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/constants"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/controllers"
	"git.pride.improwised.dev/Onboarding-2025/Krupanshi-Vaishnav/go-api/middlewares"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) error {
	mu.Lock()

	app.Use(middlewares.LogHandler(logger))

	app.Use(swagger.New(swagger.Config{
		FilePath: "./assets/swagger.json",
		Title:    "Swagger API Docs",
	}))

	err := healthCheckController(app, goqu, logger)
	if err != nil {
		return err
	}

	err = setupMoviesController(app, goqu, logger)
	if err != nil {
		return err
	}

	err = setupRatingsController(app, goqu, logger)
	if err != nil {
		return err
	}

	err = setupCastController(app, goqu, logger)
	if err != nil {
		return err
	}

	err = setupCrewController(app, goqu, logger)
	if err != nil {
		return err
	}

	mu.Unlock()
	return nil
}

func healthCheckController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	healthController, err := controllers.NewHealthController(goqu, logger)
	if err != nil {
		return err
	}

	healthz := app.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/db", healthController.Db)
	return nil
}

func setupMoviesController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	movieController, err := controllers.NewMovieController(goqu, logger)
	if err != nil {
		logger.Error("Failed to initialize MovieController", zap.Error(err))
		return err
	}

	// Register movie routes
	movieRouter := app.Group("/movies")
	movieRouter.Get("/", movieController.ListMovies)
	movieRouter.Get(fmt.Sprintf("/:%s", constants.ParamMid), movieController.GetMovieByID)
	movieRouter.Delete(fmt.Sprintf("/:%s", constants.ParamMid), movieController.DeleteMovieById)
	movieRouter.Post("/", movieController.AddMovie)
	movieRouter.Put(fmt.Sprintf("/:%s", constants.ParamMid), movieController.UpdateMovie)

	return nil

}

func setupRatingsController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	ratingController, err := controllers.NewRatingsController(goqu, logger)
	if err != nil {
		logger.Error("Failed to intialize RatingController", zap.Error(err))
		return err
	}

	ratingRouter := app.Group("ratings")

	ratingRouter.Get("/movies", ratingController.ListAllMovieRatings)
	app.Get(fmt.Sprintf("movies/:%s/ratings", constants.ParamMid), ratingController.GetRatingByMovieId)
	ratingRouter.Post(fmt.Sprintf("/user/:%s/ratings", constants.UserId), ratingController.AddRating)
	ratingRouter.Put(fmt.Sprintf("/movies/:%s/user/:%s/ratings", constants.ParamMid, constants.UserId), ratingController.UpdateRating)
	ratingRouter.Delete(fmt.Sprintf("/movies/:%s/user/:%s/ratings", constants.ParamMid, constants.UserId), ratingController.DeleteRating)

	return nil
}

func setupCastController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	castController, err := controllers.NewCastController(goqu, logger)
	if err != nil {
		logger.Error("Failed to intialize CastController", zap.Error(err))
		return err
	}

	app.Get(fmt.Sprintf("/movies/:%s/casts", constants.ParamMid), castController.ListCastMembers)
	app.Get(fmt.Sprintf("/actor/:%s/movies", constants.CastId), castController.ListMoviesByCastId)
	app.Post(fmt.Sprintf("/movies/:%s/credit/:%s/cast", constants.ParamMid, constants.CreditId), castController.AddMovieCastMember)

	return nil
}

func setupCrewController(app *fiber.App, goqu *goqu.Database, logger *zap.Logger) error {
	crewController, err := controllers.NewCrewController(goqu, logger)
	if err != nil {
		logger.Error("Failed to intialize CrewController", zap.Error(err))
		return err
	}

	app.Get(fmt.Sprintf("/movies/:%s/crew", constants.ParamMid), crewController.ListCrewMembers)
	app.Post(fmt.Sprintf("/movies/:%s/credit/:%s/crew", constants.ParamMid, constants.CreditId), crewController.AddMovieCrewMember)

	return nil
}
