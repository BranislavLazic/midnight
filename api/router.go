package api

import (
	"embed"
	"fmt"
	"github.com/branislavlazic/midnight/api/middleware"
	sess "github.com/branislavlazic/midnight/api/session"
	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/config"
	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"io/fs"
	"net/http"
	"time"
)

const sessionStoreTableName = "sessions"

type ServerSettings struct {
	Config        *config.AppConfig
	Cache         cache.Internal
	ServiceRepo   model.ServiceRepository
	UserRepo      model.UserRepository
	SessionStore  *session.Store
	TaskScheduler *task.Scheduler
	IndexFile     embed.FS
	StaticFiles   embed.FS
}

func StartServer(settings ServerSettings) error {
	secureSessionExpiry := 7 * 24 * time.Hour
	cookieKeyLookup := "cookie:" + sess.SecureCookieName
	pgStore := postgres.New(postgres.Config{
		Host:     settings.Config.DbHost,
		Port:     settings.Config.DbPort,
		Username: settings.Config.DbUser,
		Password: settings.Config.DbPassword,
		Database: settings.Config.DbName,
		Table:    sessionStoreTableName,
	})
	sessionStore := session.New(session.Config{
		Storage:        pgStore,
		Expiration:     secureSessionExpiry,
		CookieHTTPOnly: true,
		CookieSecure:   true,
		KeyLookup:      cookieKeyLookup,
		KeyGenerator:   generateSessionIdFn(settings.Config.SessionSecret),
	})

	auth := middleware.NewAuthenticator(sessionStore, settings.Config.SessionSecret)

	serviceStatusRoutes := NewServiceStatusRoutes(settings.Cache)
	serviceRoutes := NewServiceRoutes(settings.ServiceRepo, settings.TaskScheduler)
	userRoutes := NewUserRoutes(settings.UserRepo, sessionStore)

	app := fiber.New()
	// API routes
	app.Post("/v1/login", userRoutes.Login)
	app.Get("/v1/status", serviceStatusRoutes.GetStatus)
	// Authenticated routes
	app.Post("/v1/logout", userRoutes.Logout)
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
	err := app.Listen(fmt.Sprintf(":%d", settings.Config.AppPort))
	if err != nil {
		return err
	}
	return nil
}

func generateSessionIdFn(secret string) func() string {
	return func() string {
		id, err := sess.GenerateSessionID(secret)
		if err != nil {
			log.Error().Err(err).Msg("failed to generate a session id")
		}
		return id
	}
}
