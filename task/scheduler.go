package task

import (
	"github.com/branislavlazic/midnight/model"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

func InitScheduler(scheduler *gocron.Scheduler, taskProvider *Provider, serviceRepo model.ServiceRepository) error {
	services, err := serviceRepo.GetAll()
	if err != nil {
		return err
	}
	for _, s := range services {
		log.Info().Msgf(
			"initializing service %s at %s - check every %d seconds", s.Name, s.URL, s.CheckIntervalSeconds,
		)
		taskConfig := Config{ID: int64(s.ID), Name: s.Name, URL: s.URL, Timeout: s.CheckIntervalSeconds}
		_, err := scheduler.Every(s.CheckIntervalSeconds).
			Seconds().
			Do(taskProvider.NewTask(taskConfig))
		if err != nil {
			return err
		}
	}
	scheduler.StartAsync()
	return nil
}
