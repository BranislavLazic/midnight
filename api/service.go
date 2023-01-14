package api

import (
	"fmt"
	"github.com/branislavlazic/midnight/api/validation"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net/http"
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
// @Param createServiceRequest body model.CreateServiceRequest true "Create service request body"
// @Failure 400,404,409,422,500
// @Success 201
// @Router /v1/services [post]
func (lr *ServiceRoutes) CreateService(ctx *fiber.Ctx) error {
	var createServiceRequest *model.CreateServiceRequest
	if err := ctx.BodyParser(&createServiceRequest); err != nil {
		log.Debug().Err(err).Msg("failed to parse the request as service")
		return ctx.SendStatus(http.StatusBadRequest)
	}
	createServiceRequest.Sanitize()
	if err := validator.New().Struct(createServiceRequest); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(validation.ToValidationErrors(err.(validator.ValidationErrors)))
	}
	exists := lr.serviceRepo.ExistsByURL(createServiceRequest.URL)
	if exists {
		log.Debug().Msgf("a service for url %s is already registered", createServiceRequest.URL)
		return ctx.
			Status(http.StatusConflict).
			JSON(map[string]string{"error": fmt.Sprintf("A service for url %s is already registered", createServiceRequest.URL)})
	}
	ID, err := lr.serviceRepo.Create(createServiceRequest.ToPersistentService())
	if err != nil {
		log.Error().Err(err).Msg("failed to create a service")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	taskConfig := task.Config{ID: int64(ID), Name: createServiceRequest.Name, URL: createServiceRequest.URL, Timeout: createServiceRequest.CheckIntervalSeconds}
	err = lr.taskScheduler.Update(taskConfig, createServiceRequest.CheckIntervalSeconds)
	if err != nil {
		log.Error().Err(err).Msg("failed to update the task scheduler")
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	ctx.Set("Location", string(ctx.Request().Host())+ctx.Route().Path+"/"+fmt.Sprintf("%d", ID))
	return ctx.SendStatus(http.StatusCreated)
}
