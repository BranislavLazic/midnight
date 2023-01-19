package testapi

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api"
	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/config"
	"github.com/branislavlazic/midnight/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB
var Cache cache.Internal

func init() {
	cfg := TestConfig()
	DB, _ = db.GetPGDBPool(cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	Cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
}

func TestConfig() *config.AppConfig {
	return &config.AppConfig{
		DbHost:        "localhost",
		DbPort:        5432,
		DbUser:        "postgres",
		DbPassword:    "postgres",
		DbName:        "midnight",
		EnableSwagger: false,
	}
}

func InitTestApp() *fiber.App {
	serverSettings := api.ServerSettings{Config: TestConfig(), DB: DB, Cache: Cache}
	return api.InitApp(serverSettings)
}
