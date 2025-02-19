package api

import (
	"net/http"

	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type EnvironmentRoutes struct {
	repo *postgres.Repository
}

func NewEnvironmentRoutes(repo *postgres.Repository) *EnvironmentRoutes {
	return &EnvironmentRoutes{repo: repo}
}

// CreateEnvironment godoc
// @Summary Create an environment
// @Param environmentRequest body model.EnvironmentRequest true "Environment request body"
// @Failure 400,422,500
// @Success 201
// @Router /v1/environments [post]
func (er *EnvironmentRoutes) CreateEnvironment(ctx echo.Context) error {
	var envRequest *model.EnvironmentRequest
	if err := ctx.Bind(&envRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as environment")
		return ctx.NoContent(http.StatusBadRequest)
	}
	envRequest.Sanitize()
	if err := validator.New().Struct(envRequest); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	_, err := er.repo.CreateEnvironment(envRequest.ToPersistentEnvironment())
	if err != nil {
		log.Error().Err(err).Msg("failed to create a service")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusCreated)
}

// GetAllEnvironments godoc
// @Summary Get all environments
// @Failure 404
// @Success 200
// @Router /v1/environments [get]
func (er *EnvironmentRoutes) GetAllEnvironments(ctx echo.Context) error {
	envs, err := er.repo.GetAllEnvironments()
	if err != nil {
		log.Debug().Err(err).Msg("failed to fetch all environments")
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, envs)
}
