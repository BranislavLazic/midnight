package task

import (
	"github.com/branislavlazic/midnight/service"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func InitScheduler(scheduler *gocron.Scheduler, taskProvider *Provider, db *gorm.DB) error {
	serviceRepo := service.NewPgServiceRepository(db)
	services, err := serviceRepo.GetAll()
	if err != nil {
		return err
	}
	for _, s := range services {
		log.Info().Msgf(
			"initializing service %s at %s - check every %d seconds", s.Name, s.URL, s.CheckIntervalSeconds,
		)
		taskConfig := TaskConfig{ID: int64(s.ID), Name: s.Name, URL: s.URL, Timeout: s.CheckIntervalSeconds}
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
