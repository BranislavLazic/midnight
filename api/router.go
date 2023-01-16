package api

import (
	"embed"
	"fmt"
	"github.com/branislavlazic/midnight/api/middleware"
	sess "github.com/branislavlazic/midnight/api/session"
	"github.com/branislavlazic/midnight/cache"
	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"io/fs"
	"net/http"
	"time"
)

type ServerSettings struct {
	Port          int
	Cache         cache.Internal
	ServiceRepo   model.ServiceRepository
	UserRepo      model.UserRepository
	TaskScheduler *task.Scheduler
	SessionSecret string
	IndexFile     embed.FS
	StaticFiles   embed.FS
}

func StartServer(settings ServerSettings) error {
	sessionStore := session.New(session.Config{
		Expiration:     7 * 24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSecure:   true,
		KeyLookup:      "cookie:" + sess.SecureCookieName,
		KeyGenerator: func() string {
			id, _ := sess.GenerateSessionID(settings.SessionSecret)
			return id
		},
	})

	auth := middleware.NewAuthenticator(sessionStore, settings.SessionSecret)

	serviceStatusRoutes := NewServiceStatusRoutes(settings.Cache)
	serviceRoutes := NewServiceRoutes(settings.ServiceRepo, settings.TaskScheduler)
	userRoutes := NewUserRoutes(settings.UserRepo, sessionStore)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length",
	}))

	// API routes
	app.Post("/v1/login", userRoutes.Login)
	app.Get("/v1/status", serviceStatusRoutes.GetStatus)
	// Authenticated routes
	app.Get("/v1/services", auth.Authenticated(serviceRoutes.GetAllServices))
	app.Get("/v1/services/:id", auth.Authenticated(serviceRoutes.GetById))
	app.Post("/v1/services", auth.Authenticated(serviceRoutes.CreateService))
	app.Put("/v1/services/:id", auth.Authenticated(serviceRoutes.UpdateService))
	app.Delete("/v1/services/:id", auth.Authenticated(serviceRoutes.DeleteById))

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
