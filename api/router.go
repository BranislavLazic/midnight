package api

import (
	"embed"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"io/fs"
	"net/http"
)

func InitRouter(port int, cache *bigcache.BigCache, staticFiles embed.FS) error {
	app := fiber.New()
	serviceStatusRoutes := NewServiceStatusRoutes(cache)

	initStaticRoutes(app, staticFiles)

	app.Get("/v1/status", serviceStatusRoutes.GetStatus)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}

func initStaticRoutes(app *fiber.App, staticFiles embed.FS) {
	dist, _ := fs.Sub(staticFiles, "webapp/dist")
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(dist),
	}))
}
