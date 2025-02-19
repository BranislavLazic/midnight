package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ServiceRoutes struct {
	repo          *postgres.Repository
	taskScheduler *task.Scheduler
}

func NewServiceRoutes(repo *postgres.Repository, taskScheduler *task.Scheduler) *ServiceRoutes {
	return &ServiceRoutes{repo: repo, taskScheduler: taskScheduler}
}

// CreateService godoc
// @Summary Create a service
// @Param serviceRequest body model.ServiceRequest true "Service request body"
// @Failure 400,404,409,422,500
// @Success 201
// @Router /v1/services [post]
func (sr *ServiceRoutes) CreateService(ctx echo.Context) error {
	var serviceRequest *model.ServiceRequest
	if err := ctx.Bind(&serviceRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as service")
		return ctx.NoContent(http.StatusBadRequest)
	}
	serviceRequest.Sanitize()
	if err := validator.New().Struct(serviceRequest); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	exists := sr.repo.ServiceExistsByURL(serviceRequest.URL)
	if exists {
		log.Debug().Msgf("a service for url %s is already registered", serviceRequest.URL)
		return ctx.
			JSON(http.StatusConflict, map[string]string{"error": fmt.Sprintf("A service for url %s is already registered", serviceRequest.URL)})
	}
	env, err := sr.repo.GetEnvironmentByID(model.EnvironmentID(serviceRequest.EnvironmentID))
	if err != nil {
		log.Warn().Err(err).Msg("environment will not be set")
	}
	ID, err := sr.repo.CreateService(serviceRequest.ToPersistentService(0, env))

	if err != nil {
		log.Error().Err(err).Msg("failed to create a service")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	taskConfig := task.Config{
		ID:           int64(ID),
		Name:         serviceRequest.Name,
		URL:          serviceRequest.URL,
		ResponseBody: serviceRequest.ResponseBody,
		Timeout:      serviceRequest.CheckIntervalSeconds,
	}
	err = sr.taskScheduler.Add(taskConfig, serviceRequest.CheckIntervalSeconds)
	if err != nil {
		log.Error().Err(err).Msg("failed to add the task")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusCreated)
}

// UpdateService godoc
// @Summary Update a service
// @Param serviceRequest body model.ServiceRequest true "Service request body"
// @Failure 400,404,422,500
// @Success 200
// @Router /v1/services/{id} [post]
func (sr *ServiceRoutes) UpdateService(ctx echo.Context) error {
	idStr := ctx.Param("id")
	ID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}
	var serviceRequest *model.ServiceRequest
	if err := ctx.Bind(&serviceRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as service")
		return ctx.NoContent(http.StatusBadRequest)
	}
	serviceRequest.Sanitize()
	if err := validator.New().Struct(serviceRequest); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	env, err := sr.repo.GetEnvironmentByID(model.EnvironmentID(serviceRequest.EnvironmentID))
	if err != nil {
		log.Warn().Err(err).Msg("environment will not be set")
	}
	err = sr.repo.UpdateService(serviceRequest.ToPersistentService(model.ServiceID(ID), env))
	if err != nil {
		log.Error().Err(err).Msg("failed to update the service")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	taskConfig := task.Config{
		ID:           int64(ID),
		Name:         serviceRequest.Name,
		URL:          serviceRequest.URL,
		ResponseBody: serviceRequest.ResponseBody,
		Timeout:      serviceRequest.CheckIntervalSeconds,
	}
	err = sr.taskScheduler.Update(taskConfig, serviceRequest.CheckIntervalSeconds)
	if err != nil {
		log.Error().Err(err).Msg("failed to update the task scheduler")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

// GetAllServices godoc
// @Summary Get all services
// @Failure 404
// @Success 200
// @Router /v1/services [get]
func (sr *ServiceRoutes) GetAllServices(ctx echo.Context) error {
	services, err := sr.repo.GetAllServices()
	if err != nil {
		log.Debug().Err(err).Msg("failed to fetch all services")
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, services)
}

// GetById godoc
// @Summary Get a service by id
// @Param id path string true "Service ID"
// @Failure 404
// @Success 200
// @Router /v1/services/{id} [get]
func (sr *ServiceRoutes) GetById(ctx echo.Context) error {
	idStr := ctx.Param("id")
	ID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}
	service, err := sr.repo.GetServiceById(model.ServiceID(ID))
	if err != nil {
		log.Debug().Err(err).Msgf("failed to find a service by id %d", ID)
		return ctx.NoContent(http.StatusNotFound)
	}
	return ctx.JSON(http.StatusOK, service)
}

// DeleteById godoc
// @Summary Delete a service
// @Produce json
// @Param id path string true "Service ID"
// @Failure 400,404
// @Success 204
// @Router /v1/services/{id} [delete]
func (sr *ServiceRoutes) DeleteById(ctx echo.Context) error {
	idStr := ctx.Param("id")
	ID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.NoContent(http.StatusBadRequest)
		return err
	}
	err = sr.repo.DeleteServiceById(model.ServiceID(ID))
	if err != nil {
		log.Debug().Err(err).Msg("failed to delete the service. not found.")
		return ctx.NoContent(http.StatusNotFound)
	}
	err = sr.taskScheduler.Remove(int64(ID))
	if err != nil {
		log.Debug().Err(err).Msgf("failed to remove task with id %d.", ID)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusNoContent)
}
