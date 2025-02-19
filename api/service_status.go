package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/task"
	"github.com/labstack/echo/v4"
)

type ServiceStatusRoutes struct {
	cache cache.Internal
}

func NewServiceStatusRoutes(cache cache.Internal) *ServiceStatusRoutes {
	return &ServiceStatusRoutes{cache: cache}
}

// GetStatus godoc
// @Summary Get status
// @Failure 404
// @Success 200
// @Router /v1/status [get]
func (lr *ServiceStatusRoutes) GetStatus(ctx echo.Context) error {
	bytes, err := lr.cache.Get(task.ServiceStatusCacheName)
	if err != nil {
		return ctx.JSON(http.StatusOK, []task.ServiceStatus{})
	}
	serviceStatuses, err := sortServiceStatuses(bytes)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, serviceStatuses)
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
