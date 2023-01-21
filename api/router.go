package api

import (
	"embed"
	"fmt"
	"github.com/branislavlazic/midnight/api/middleware"
	sess "github.com/branislavlazic/midnight/api/session"
	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/config"
	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/session"
	pg "github.com/gofiber/storage/postgres"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
	"io/fs"
	"net/http"
	"time"
)

const sessionStoreTableName = "sessions"

type ServerSettings struct {
	Config      *config.AppConfig
	DB          *gorm.DB
	Cache       cache.Internal
	IndexFile   embed.FS
	StaticFiles embed.FS
}

func InitApp(settings ServerSettings) *fiber.App {
	serviceRepo := postgres.NewServiceRepository(settings.DB)
	userRepo := postgres.NewUserRepository(settings.DB)
	envRepo := postgres.NewEnvironmentRepository(settings.DB)

	scheduler := gocron.NewScheduler(time.UTC)
	taskProvider := task.NewProvider(settings.Cache)
	taskScheduler := task.NewScheduler(scheduler, taskProvider, serviceRepo)
	err := taskScheduler.RunAll()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize task scheduler")
	}

	secureSessionExpiry := 7 * 24 * time.Hour
	cookieKeyLookup := "cookie:" + sess.SecureCookieName
	pgStore := pg.New(pg.Config{
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

	authFn := middleware.Authenticated(sessionStore, settings.Config.SessionSecret)

	authRoutes := NewAuthRoutes(userRepo, sessionStore)
	serviceStatusRoutes := NewServiceStatusRoutes(settings.Cache)
	serviceRoutes := NewServiceRoutes(serviceRepo, taskScheduler)
	envRoutes := NewEnvironmentRoutes(envRepo)

	app := fiber.New()
	// API routes
	app.Post("/v1/login", authRoutes.Login)
	app.Get("/v1/status", serviceStatusRoutes.GetStatus)
	// Authenticated routes
	app.Post("/v1/logout", authRoutes.Logout)
	app.Get("/v1/services", serviceRoutes.GetAllServices)
	app.Get("/v1/services/:id", authFn, serviceRoutes.GetById)
	app.Post("/v1/services", authFn, serviceRoutes.CreateService)
	app.Put("/v1/services/:id", authFn, serviceRoutes.UpdateService)
	app.Delete("/v1/services/:id", authFn, serviceRoutes.DeleteById)
	app.Post("/v1/environments", authFn, envRoutes.CreateEnvironment)
	app.Get("/v1/environments", envRoutes.GetAllEnvironments)

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
	return app
}

func StartServer(settings ServerSettings) error {
	err := InitApp(settings).Listen(fmt.Sprintf(":%d", settings.Config.AppPort))
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
