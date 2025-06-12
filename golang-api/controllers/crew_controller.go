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

type CrewController struct {
	crewModel *models.CrewModel
	logger    *zap.Logger
}

func NewCrewController(logger *zap.Logger) (*CrewController, error) {
	model := models.NewCrewModel()
	return &CrewController{
		crewModel: model,
		logger:    logger,
	}, nil
}

// ListCrewMembers retrieves all crew members by movieID
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
//	500: GenericErrorResponse
func (ctrl *CrewController) ListCrewMembers(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)
	crewMembers, err := ctrl.crewModel.ListCrewMembers(movieId)
	if err != nil {
		ctrl.logger.Error(constants.LoadCreditsError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.LoadCreditsError)
	}

	return utils.JSONSuccess(c, http.StatusOK, crewMembers)

}

// UpdateCrewMember updates a crew member details
// swagger:route PUT /movies/{movieId}/crew/{crewId} Crew UpdateCrewMember
//
// Updates an existing crew member.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
// - RequestUpdateCrewMember
//
// Responses:
//
//	200: GenericSuccessResponse
//	400: ValidationErrorResponse
//	500: GenericErrorResponse
func (ctrl *CrewController) UpdateCrewMember(c *fiber.Ctx) error {
	movieId := c.Params(constants.MovieId)
	crewId := c.Params(constants.CrewId)

	var validate = validator.New()
	var crew models.CrewMember

	// Unmarshal JSON data into the movie struct
	if err := json.Unmarshal(c.Body(), &crew); err != nil {
		ctrl.logger.Error(constants.InvalidRequestBody, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidRequestBody)
	}

	// Validate the movie struct using validator
	if err := validate.Struct(crew); err != nil {
		ctrl.logger.Error(constants.ValidationFailed, zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	if err := ctrl.crewModel.UpdateCrewMember(movieId, crewId, crew); err != nil {
		ctrl.logger.Error(constants.UpdateCrewError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.UpdateCrewError)
	}

	return utils.JSONSuccess(c, http.StatusOK, constants.UpdateCrewSuccess)
}
