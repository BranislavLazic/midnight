package main

import (
	"context"
	"embed"
	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api"
	"github.com/branislavlazic/midnight/config"
	"github.com/branislavlazic/midnight/db"
	"github.com/rs/zerolog/log"
	"time"
)

//go:embed webapp/dist
var uiStaticFiles embed.FS

//go:embed webapp/dist/index.html
var indexFile embed.FS

//go:embed migrations/*.sql
var dbMigrations embed.FS

func main() {
	cfg := config.GetAppConfig()
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize cache")
	}
	pgDb, err := db.GetPGDBPool(cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}
	err = db.RunMigrations(pgDb, dbMigrations)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to run the migrations")
	}

	serverSettings := api.ServerSettings{
		Config:      cfg,
		DB:          pgDb,
		Cache:       cache,
		IndexFile:   indexFile,
		StaticFiles: uiStaticFiles,
	}
	err = api.StartServer(serverSettings)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
	}
}
