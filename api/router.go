package api

import (
	"embed"
	"fmt"
	"github.com/branislavlazic/midnight/cache"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"io/fs"
	"net/http"

	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type ServerSettings struct {
	Port          int
	Cache         cache.Internal
	ServiceRepo   model.ServiceRepository
	TaskScheduler *task.Scheduler
	IndexFile     embed.FS
	StaticFiles   embed.FS
}

func StartServer(settings ServerSettings) error {
	serviceStatusRoutes := NewServiceStatusRoutes(settings.Cache)
	serviceRoutes := NewServiceRoutes(settings.ServiceRepo, settings.TaskScheduler)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length",
	}))

	// API routes
	app.Get("/v1/status", serviceStatusRoutes.GetStatus)
	app.Get("/v1/services", serviceRoutes.GetAllServices)
	app.Get("/v1/services/:id", serviceRoutes.GetById)
	app.Post("/v1/services", serviceRoutes.CreateService)
	app.Put("/v1/services/:id", serviceRoutes.UpdateService)
	app.Delete("/v1/services/:id", serviceRoutes.DeleteById)

	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Webapp static files
	assetsSubDir, _ := fs.Sub(settings.StaticFiles, "webapp/dist")
	indexSubDir, _ := fs.Sub(settings.IndexFile, "webapp/dist")
	assetsFileDir := filesystem.New(filesystem.Config{
		Root: http.FS(assetsSubDir),
	})
	// Serve assets directory
	app.Get("/assets/*", assetsFileDir)
	// Serve svg logo
	app.Get("/*.svg", assetsFileDir)
	// Serve index.html
	app.Get("/*", func(ctx *fiber.Ctx) error {
		return filesystem.SendFile(ctx, http.FS(indexSubDir), "/index.html")
	})

	// Start server
	err := app.Listen(fmt.Sprintf(":%d", settings.Port))
	if err != nil {
		return err
	}
	return nil
}
