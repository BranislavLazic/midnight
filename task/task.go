package task

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

const LivenessCacheName = "liveness"

type Liveness struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Version    string `json:"version,omitempty"`
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
}

type LivenessResponse struct {
	Version string `json:"version"`
}

type TaskConfig struct {
	ID      int64
	Name    string
	URL     string
	Timeout int
}

type Provider struct {
	cache *bigcache.BigCache
}

func NewProvider(cache *bigcache.BigCache) *Provider {
	return &Provider{cache: cache}
}

func (tp *Provider) NewTask(config TaskConfig) func() {
	return func() {
		req, err := http.NewRequest(http.MethodGet, config.URL, nil)
		if err != nil {
			log.Warn().Err(err).Msg("could not create the request")
		}
		client := http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		}
		res, err := client.Do(req)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get a response")
			tp.handleFailureResponse(config.ID, config.Name, config.URL)
		} else {
			tp.handleSuccessResponse(res, config.ID, config.Name, config.URL)
		}
	}
}

func (tp *Provider) handleSuccessResponse(res *http.Response, id int64, name, url string) {
	log.Debug().Msg(res.Status)
	var livenessResponse LivenessResponse
	err := json.NewDecoder(res.Body).Decode(&livenessResponse)
	if err != nil {
		log.Warn().Err(err).Msg("failed to extract request body")
	}
	err = tp.setLiveness(
		Liveness{ID: id, Name: name, URL: url, Version: livenessResponse.Version, Status: res.Status, StatusCode: res.StatusCode},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to set task")
	}
}

func (tp *Provider) handleFailureResponse(id int64, name, url string) {
	err := tp.setLiveness(Liveness{ID: id, Name: name, URL: url, Status: "404 Not Found", StatusCode: 404})
	if err != nil {
		log.Error().Err(err).Msg("failed to set task")
	}
}

func (tp *Provider) setLiveness(liveness Liveness) error {
	bytes, err := tp.cache.Get(LivenessCacheName)
	if err != nil {
		return tp.cache.Set(LivenessCacheName, serializeLiveness(map[int64]Liveness{liveness.ID: liveness}))
	} else {
		var livenesses map[int64]Liveness
		err = json.Unmarshal(bytes, &livenesses)
		if err != nil {
			return err
		}
		livenesses[liveness.ID] = liveness
		return tp.cache.Set(LivenessCacheName, serializeLiveness(livenesses))
	}
}

func serializeLiveness(liveness map[int64]Liveness) []byte {
	bytes, _ := json.Marshal(&liveness)
	return bytes
}
