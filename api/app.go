package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/config"
	_ "github.com/branislavlazic/midnight/docs"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-co-op/gocron"

	echoSwagger "github.com/swaggo/echo-swagger"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ServerSettings struct {
	Config *config.AppConfig
	DB     *gorm.DB
	Cache  cache.Internal
}

type App struct {
	repository *postgres.Repository
	settings   ServerSettings
}

func NewApp(repository *postgres.Repository, settings ServerSettings) *App {
	return &App{repository: repository, settings: settings}
}

func (a *App) InitApi() *echo.Echo {
	scheduler := gocron.NewScheduler(time.UTC)
	taskProvider := task.NewProvider(a.settings.Cache)
	taskScheduler := task.NewScheduler(scheduler, taskProvider, a.repository)
	err := taskScheduler.RunAll()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize task scheduler")
	}

	authRoutes := NewAuthRoutes(a.repository, a.settings.Config)
	serviceStatusRoutes := NewServiceStatusRoutes(a.settings.Cache)
	serviceRoutes := NewServiceRoutes(a.repository, taskScheduler)
	envRoutes := NewEnvironmentRoutes(a.repository)

	e := echo.New()

	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"*"},
	}))
	// API routes
	insecure := e.Group("/v1")
	insecure.POST("/login", authRoutes.Login)
	insecure.GET("/status", serviceStatusRoutes.GetStatus)
	insecure.GET("/environments", envRoutes.GetAllEnvironments)

	// Authenticated routes
	secure := e.Group("/v1")
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(a.settings.Config.SessionSecret),
	}
	secure.Use(echojwt.WithConfig(jwtConfig))
	secure.GET("/services", serviceRoutes.GetAllServices)
	secure.GET("/services/:id", serviceRoutes.GetById)
	secure.POST("/services", serviceRoutes.CreateService)
	secure.PUT("/services/:id", serviceRoutes.UpdateService)
	secure.DELETE("/services/:id", serviceRoutes.DeleteById)
	secure.POST("/environments", envRoutes.CreateEnvironment)

	// Swagger
	swagger := e.Group("/swagger")
	swagger.GET("/*", echoSwagger.WrapHandler)

	e.File("/*", "webapp/dist/index.html")
	e.Use(echoMw.StaticWithConfig(echoMw.StaticConfig{
		Root:   "webapp/dist",
		Browse: true,
	}))

	return e
}
