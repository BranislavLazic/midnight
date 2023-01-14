package main

import (
	"context"
	"embed"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

//go:embed webapp/dist
var embedDirStatic embed.FS

func main() {
	scheduler := gocron.NewScheduler(time.UTC)
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))

	taskProvider := task.NewProvider(cache)
	task1 := taskProvider.NewTask("http://localhost:8000", 5)
	task2 := taskProvider.NewTask("https://google.rs", 5)

	scheduler.Every(5).Seconds().Do(task1)
	scheduler.Every(5).Seconds().Do(task2)

	scheduler.StartAsync()
	err := api.InitRouter(8000, cache, embedDirStatic)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
	}
}
