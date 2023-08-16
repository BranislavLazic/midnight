package api

import (
	"fmt"
	"net/http"

	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type EnvironmentRoutes struct {
	envRepo model.EnvironmentRepository
}

func NewEnvironmentRoutes(envRepo model.EnvironmentRepository) *EnvironmentRoutes {
	return &EnvironmentRoutes{envRepo: envRepo}
}

// CreateEnvironment godoc
// @Summary Create an environment
// @Param environmentRequest body model.EnvironmentRequest true "Environment request body"
// @Failure 400,422,500
// @Success 201
// @Router /v1/environments [post]
func (er *EnvironmentRoutes) CreateEnvironment(ctx *fiber.Ctx) error {
	var envRequest *model.EnvironmentRequest
	if err := ctx.BodyParser(&envRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as environment")
		return ctx.SendStatus(http.StatusBadRequest)
	}
	envRequest.Sanitize()
	if err := validator.New().Struct(envRequest); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	ID, err := er.envRepo.Create(envRequest.ToPersistentEnvironment())
	if err != nil {
		log.Error().Err(err).Msg("failed to create a service")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	ctx.Set("Location", string(ctx.Request().Host())+ctx.Route().Path+"/"+fmt.Sprintf("%d", ID))
	return ctx.SendStatus(http.StatusCreated)
}

// GetAllEnvironments godoc
// @Summary Get all environments
// @Failure 404
// @Success 200
// @Router /v1/environments [get]
func (er *EnvironmentRoutes) GetAllEnvironments(ctx *fiber.Ctx) error {
	envs, err := er.envRepo.GetAll()
	if err != nil {
		log.Debug().Err(err).Msg("failed to fetch all environments")
		return ctx.SendStatus(http.StatusNotFound)
	}
	return ctx.Status(http.StatusOK).JSON(envs)
}
