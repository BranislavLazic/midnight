package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/allegro/bigcache/v3"
	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func InitRouter(port int, cache *bigcache.BigCache, serviceRepo model.ServiceRepository, taskScheduler *task.Scheduler, staticFiles embed.FS) error {
	serviceStatusRoutes := NewServiceStatusRoutes(cache)
	serviceRoutes := NewServiceRoutes(serviceRepo, taskScheduler)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length",
	}))
	dist, _ := fs.Sub(staticFiles, "webapp/dist")
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(dist),
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/v1/status", serviceStatusRoutes.GetStatus)
	app.Post("/v1/services", serviceRoutes.CreateService)
	app.Get("/v1/services", serviceRoutes.GetAllServices)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}
