package api

import (
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sort"
)

type LivenessRoutes struct {
	cache *bigcache.BigCache
}

func NewLivenessRoutes(cache *bigcache.BigCache) *LivenessRoutes {
	return &LivenessRoutes{cache: cache}
}

func (lr *LivenessRoutes) GetStatus(ctx *fiber.Ctx) error {
	bytes, err := lr.cache.Get(task.LivenessCacheName)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	livenesses, err := getSortedLivenesses(bytes)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.Status(http.StatusOK).JSON(livenesses)
}

func getSortedLivenesses(bytes []byte) ([]task.Liveness, error) {
	var livenessesMap map[string]task.Liveness
	err := json.Unmarshal(bytes, &livenessesMap)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(livenessesMap))
	for k := range livenessesMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var livenesses []task.Liveness
	for _, key := range keys {
		livenesses = append(livenesses, livenessesMap[key])
	}
	return livenesses, nil
}
