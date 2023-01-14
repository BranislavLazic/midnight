package main

import (
	"context"
	"embed"
	"github.com/branislavlazic/midnight/db"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

//go:embed webapp/dist
var uiStaticFiles embed.FS

//go:embed migrations/*.sql
var dbMigrations embed.FS

func main() {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize cache")
	}
	pgDb, err := db.GetPGDBPool("localhost", "postgres", "postgres", "midnight", 5432)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	err = db.RunMigrations(pgDb, dbMigrations)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to run the migrations")
	}

	scheduler := gocron.NewScheduler(time.UTC)
	taskProvider := task.NewProvider(cache)
	err = task.InitScheduler(scheduler, taskProvider, pgDb)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize scheduler")
	}

	err = api.InitRouter(8000, cache, uiStaticFiles)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
	}
}
