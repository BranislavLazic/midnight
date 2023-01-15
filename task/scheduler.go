package task

import (
	"github.com/branislavlazic/midnight/model"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"strconv"
)

type Scheduler struct {
	scheduler    *gocron.Scheduler
	taskProvider *Provider
	serviceRepo  model.ServiceRepository
}

func NewScheduler(scheduler *gocron.Scheduler, taskProvider *Provider, serviceRepo model.ServiceRepository) *Scheduler {
	return &Scheduler{scheduler: scheduler, taskProvider: taskProvider, serviceRepo: serviceRepo}
}

func (s *Scheduler) RunAll() error {
	services, err := s.serviceRepo.GetAll()
	if err != nil {
		return err
	}
	for _, service := range services {
		log.Info().Msgf(
			"initializing service %service at %service - check every %d seconds", service.Name, service.URL, service.CheckIntervalSeconds,
		)
		taskConfig := Config{ID: int64(service.ID), Name: service.Name, URL: service.URL, Timeout: service.CheckIntervalSeconds}
		_, err := s.scheduler.Every(service.CheckIntervalSeconds).
			Tag(strconv.FormatInt(int64(service.ID), 10)).
			Seconds().
			Do(s.taskProvider.NewTask(taskConfig))
		if err != nil {
			return err
		}
	}
	s.scheduler.StartAsync()
	return nil
}

func (s *Scheduler) Update(cfg Config, checkIntervalSeconds int) error {
	job, err := s.scheduler.Every(checkIntervalSeconds).
		Tag(strconv.FormatInt(cfg.ID, 10)).
		Seconds().
		Do(s.taskProvider.NewTask(cfg))
	if err != nil {
		return err
	}
	_, err = s.scheduler.Job(job).Update()
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) Remove(ID int64) error {
	_ = s.taskProvider.RemoveTask(ID)
	return s.scheduler.RemoveByTag(strconv.FormatInt(ID, 10))
}
