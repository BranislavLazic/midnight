package task

import (
	"encoding/json"
	"github.com/branislavlazic/midnight/cache"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const ServiceStatusCacheName = "service-status"

type ServiceStatus struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Version    string `json:"version,omitempty"`
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
}

type ServiceStatusResponse struct {
	Version string `json:"version"`
}

type Config struct {
	ID      int64
	Name    string
	URL     string
	Timeout int
}

type Provider struct {
	cache cache.Internal
}

func NewProvider(c cache.Internal) *Provider {
	return &Provider{cache: c}
}

// NewTask builds a function based on a configuration
// that should be run by a scheduler.
// The function is reading the configuration and
// an HTTP client is pinging a service endpoint.
// Depending on the response, the underlying cache will be updated
// with the latest service status.
func (tp *Provider) NewTask(config Config) func() {
	return func() {
		req, err := http.NewRequest(http.MethodGet, config.URL, nil)
		if err != nil {
			log.Warn().Err(err).Msg("failed to create the request")
		}
		client := http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		}
		res, err := client.Do(req)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get a response")
			err := SaveServiceStatus(
				tp.cache,
				ServiceStatus{ID: config.ID, Name: config.Name, URL: config.URL, Status: "404 Not Found", StatusCode: 404},
			)
			if err != nil {
				log.Error().Err(err).Msg("failed to set the task")
			}
		} else {
			log.Debug().Msg(res.Status)
			var serviceStatusResponse ServiceStatusResponse
			err := json.NewDecoder(res.Body).Decode(&serviceStatusResponse)
			if err != nil {
				//log.Debug().Err(err).Msg("failed to extract the request body")
			}
			err = SaveServiceStatus(
				tp.cache,
				ServiceStatus{ID: config.ID, Name: config.Name, URL: config.URL, Version: serviceStatusResponse.Version, Status: res.Status, StatusCode: res.StatusCode},
			)
			if err != nil {
				log.Error().Err(err).Msg("failed to set the task")
			}
		}
	}
}

// RemoveTask removes the task from the underlying cache
func (tp *Provider) RemoveTask(ID int64) error {
	bytes, err := tp.cache.Get(ServiceStatusCacheName)
	if err != nil {
		return nil
	}
	var serviceStatuses map[int64]ServiceStatus
	err = json.Unmarshal(bytes, &serviceStatuses)
	if err != nil {
		return err
	}
	delete(serviceStatuses, ID)
	return tp.cache.Set(ServiceStatusCacheName, serializeServiceStatus(serviceStatuses))
}

func SaveServiceStatus(cache cache.Internal, serviceStatus ServiceStatus) error {
	bytes, err := cache.Get(ServiceStatusCacheName)
	if err != nil {
		return cache.Set(
			ServiceStatusCacheName,
			serializeServiceStatus(map[int64]ServiceStatus{serviceStatus.ID: serviceStatus}),
		)
	} else {
		var serviceStatuses map[int64]ServiceStatus
		err = json.Unmarshal(bytes, &serviceStatuses)
		if err != nil {
			return err
		}
		serviceStatuses[serviceStatus.ID] = serviceStatus
		return cache.Set(ServiceStatusCacheName, serializeServiceStatus(serviceStatuses))
	}
}

func serializeServiceStatus(serviceStatuses map[int64]ServiceStatus) []byte {
	bytes, _ := json.Marshal(&serviceStatuses)
	return bytes
}
