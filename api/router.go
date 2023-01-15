package api

import (
	"embed"
	"fmt"
	"github.com/branislavlazic/midnight/api/middleware/filesys"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"io/fs"
	"net/http"

	"github.com/allegro/bigcache/v3"
	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func InitRouter(port int, cache *bigcache.BigCache, serviceRepo model.ServiceRepository, taskScheduler *task.Scheduler, indexFile embed.FS, staticFiles embed.FS) error {
	serviceStatusRoutes := NewServiceStatusRoutes(cache)
	serviceRoutes := NewServiceRoutes(serviceRepo, taskScheduler)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length",
	}))

	app.Get("/v1/status", serviceStatusRoutes.GetStatus)
	app.Post("/v1/services", serviceRoutes.CreateService)
	app.Get("/v1/services", serviceRoutes.GetAllServices)
	app.Get("/v1/services/:id", serviceRoutes.GetById)

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Webapp static files
	assetsSubDir, _ := fs.Sub(staticFiles, "webapp/dist")
	indexSubDir, _ := fs.Sub(indexFile, "webapp/dist")
	assetsFileDir := filesystem.New(filesystem.Config{
		Root: http.FS(assetsSubDir),
	})
	// Serve assets directory
	app.Get("/assets/*", assetsFileDir)
	// Serve svg logo
	app.Get("/*.svg", assetsFileDir)
	// Serve index.html
	app.Get("/*", filesys.NewFileSys(filesystem.Config{
		Root: http.FS(indexSubDir),
	}))

	// Start
	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}
