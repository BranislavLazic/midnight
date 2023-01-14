package api

import (
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2"
)

func InitRouter(port int, cache *bigcache.BigCache) error {
	app := fiber.New()
	livenessRoutes := NewLivenessRoutes(cache)

	app.Get("/", livenessRoutes.GetStatuses)

	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}
