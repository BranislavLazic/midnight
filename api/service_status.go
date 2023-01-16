package api

import (
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sort"
)

type ServiceStatusRoutes struct {
	cache *bigcache.BigCache
}

func NewServiceStatusRoutes(cache *bigcache.BigCache) *ServiceStatusRoutes {
	return &ServiceStatusRoutes{cache: cache}
}

// GetStatus godoc
// @Summary Get status
// @Failure 404
// @Success 200
// @Router /v1/status [get]
func (lr *ServiceStatusRoutes) GetStatus(ctx *fiber.Ctx) error {
	bytes, err := lr.cache.Get(task.ServiceStatusCacheName)
	if err != nil {
		return ctx.Status(http.StatusOK).JSON([]task.ServiceStatus{})
	}
	serviceStatuses, err := sortServiceStatuses(bytes)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	return ctx.Status(http.StatusOK).JSON(serviceStatuses)
}

func sortServiceStatuses(bytes []byte) ([]task.ServiceStatus, error) {
	var serviceStatusesMap map[int64]task.ServiceStatus
	err := json.Unmarshal(bytes, &serviceStatusesMap)
	if err != nil {
		return nil, err
	}
	keys := make([]int64, 0, len(serviceStatusesMap))
	for k := range serviceStatusesMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	serviceStatuses := []task.ServiceStatus{}
	for _, key := range keys {
		serviceStatuses = append(serviceStatuses, serviceStatusesMap[key])
	}
	return serviceStatuses, nil
}
