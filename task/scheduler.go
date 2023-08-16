package task

import (
	"strconv"

	"github.com/branislavlazic/midnight/model"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	scheduler    *gocron.Scheduler
	taskProvider *Provider
	serviceRepo  model.ServiceRepository
}

func NewScheduler(scheduler *gocron.Scheduler, taskProvider *Provider, serviceRepo model.ServiceRepository) *Scheduler {
	return &Scheduler{scheduler: scheduler, taskProvider: taskProvider, serviceRepo: serviceRepo}
}

// RunAll reads services settings, creates task configurations,
// spawns jobs, and runs the scheduler in the end.
func (s *Scheduler) RunAll() error {
	services, err := s.serviceRepo.GetAll()
	if err != nil {
		return err
	}
	for _, service := range services {
		log.Info().Msgf(
			"initializing service %service at %service - check every %d seconds", service.Name, service.URL, service.CheckIntervalSeconds,
		)
		taskConfig := Config{ID: int64(service.ID), Name: service.Name, URL: service.URL, Environment: service.Environment, ResponseBody: service.ResponseBody, Timeout: service.CheckIntervalSeconds}
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

// Add updates the scheduler with a new job
// that was provided through a task configuration.
func (s *Scheduler) Add(cfg Config, checkIntervalSeconds int) error {
	job, err := s.scheduler.Every(checkIntervalSeconds).
		Tag(strconv.FormatInt(cfg.ID, 10)).
		Seconds().
		Do(s.taskProvider.NewTask(cfg))
	if err != nil {
		return err
	}
	if !s.scheduler.IsRunning() {
		s.scheduler.StartAsync()
	}
	_, err = s.scheduler.Job(job).Update()
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) Update(cfg Config, checkIntervalSeconds int) error {
	err := s.Remove(cfg.ID)
	if err != nil {
		return err
	}
	return s.Add(cfg, checkIntervalSeconds)
}

// Remove removes the task from the cache by its id
// and removes its job from the scheduler.
func (s *Scheduler) Remove(ID int64) error {
	_ = s.taskProvider.RemoveTask(ID)
	return s.scheduler.RemoveByTag(strconv.FormatInt(ID, 10))
}
