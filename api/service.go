package api

import (
	"fmt"
	"net/http"

	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type ServiceRoutes struct {
	serviceRepo   model.ServiceRepository
	taskScheduler *task.Scheduler
}

func NewServiceRoutes(serviceRepo model.ServiceRepository, taskScheduler *task.Scheduler) *ServiceRoutes {
	return &ServiceRoutes{serviceRepo: serviceRepo, taskScheduler: taskScheduler}
}

// CreateService godoc
// @Summary Create a service
// @Param serviceRequest body model.ServiceRequest true "Service request body"
// @Failure 400,404,409,422,500
// @Success 201
// @Router /v1/services [post]
func (sr *ServiceRoutes) CreateService(ctx *fiber.Ctx) error {
	var serviceRequest *model.ServiceRequest
	if err := ctx.BodyParser(&serviceRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as service")
		return ctx.SendStatus(http.StatusBadRequest)
	}
	serviceRequest.Sanitize()
	if err := validator.New().Struct(serviceRequest); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	exists := sr.serviceRepo.ExistsByURL(serviceRequest.URL)
	if exists {
		log.Debug().Msgf("a service for url %s is already registered", serviceRequest.URL)
		return ctx.
			Status(http.StatusConflict).
			JSON(map[string]string{"error": fmt.Sprintf("A service for url %s is already registered", serviceRequest.URL)})
	}
	ID, err := sr.serviceRepo.Create(serviceRequest.ToPersistentService(0))
	if err != nil {
		log.Error().Err(err).Msg("failed to create a service")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	taskConfig := task.Config{ID: int64(ID), Name: serviceRequest.Name, URL: serviceRequest.URL, Timeout: serviceRequest.CheckIntervalSeconds}
	err = sr.taskScheduler.Add(taskConfig, serviceRequest.CheckIntervalSeconds)
	if err != nil {
		log.Error().Err(err).Msg("failed to add the task")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	ctx.Set("Location", string(ctx.Request().Host())+ctx.Route().Path+"/"+fmt.Sprintf("%d", ID))
	return ctx.SendStatus(http.StatusCreated)
}

// UpdateService godoc
// @Summary Update a service
// @Param serviceRequest body model.ServiceRequest true "Service request body"
// @Failure 400,404,422,500
// @Success 200
// @Router /v1/services/{id} [post]
func (sr *ServiceRoutes) UpdateService(ctx *fiber.Ctx) error {
	ID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	var serviceRequest *model.ServiceRequest
	if err := ctx.BodyParser(&serviceRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as service")
		return ctx.SendStatus(http.StatusBadRequest)
	}
	serviceRequest.Sanitize()
	if err := validator.New().Struct(serviceRequest); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	err = sr.serviceRepo.Update(serviceRequest.ToPersistentService(model.ServiceID(ID)))
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	taskConfig := task.Config{ID: int64(ID), Name: serviceRequest.Name, URL: serviceRequest.URL, Timeout: serviceRequest.CheckIntervalSeconds}
	err = sr.taskScheduler.Update(taskConfig, serviceRequest.CheckIntervalSeconds)
	if err != nil {
		log.Error().Err(err).Msg("failed to update the task scheduler")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.SendStatus(http.StatusOK)
}

// GetAllServices godoc
// @Summary Get all services
// @Failure 404
// @Success 200
// @Router /v1/services [get]
func (sr *ServiceRoutes) GetAllServices(ctx *fiber.Ctx) error {
	services, err := sr.serviceRepo.GetAll()
	if err != nil {
		log.Debug().Err(err).Msg("failed to fetch all services")
		return ctx.SendStatus(http.StatusNotFound)
	}
	return ctx.Status(http.StatusOK).JSON(services)
}

// GetById godoc
// @Summary Get a service by id
// @Param id path string true "Service ID"
// @Failure 404
// @Success 200
// @Router /v1/services/{id} [get]
func (sr *ServiceRoutes) GetById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	service, err := sr.serviceRepo.GetById(model.ServiceID(id))
	if err != nil {
		log.Debug().Err(err).Msgf("failed to find a service by id %d", id)
		return ctx.SendStatus(http.StatusNotFound)
	}
	return ctx.Status(http.StatusOK).JSON(service)
}

// DeleteById godoc
// @Summary Delete a service
// @Produce json
// @Param id path string true "Service ID"
// @Failure 400,404
// @Success 204
// @Router /v1/services/{id} [delete]
func (sr *ServiceRoutes) DeleteById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	err = sr.serviceRepo.DeleteById(model.ServiceID(id))
	if err != nil {
		log.Debug().Err(err).Msg("failed to delete the service. not found.")
		return ctx.SendStatus(http.StatusNotFound)
	}
	err = sr.taskScheduler.Remove(int64(id))
	if err != nil {
		log.Debug().Err(err).Msgf("failed to remove task with id %d.", id)
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.SendStatus(http.StatusNoContent)
}
