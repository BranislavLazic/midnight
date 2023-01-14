package task

import (
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const LivenessCacheName = "liveness"

type Liveness struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

type Provider struct {
	cache *bigcache.BigCache
}

func NewProvider(cache *bigcache.BigCache) *Provider {
	return &Provider{cache: cache}
}

func (tp *Provider) NewTask(url string, everySeconds int) func() {
	return func() {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Warn().Err(err).Msg("could not create the request")
		}
		client := http.Client{
			Timeout: time.Duration(everySeconds) * time.Second,
		}
		res, err := client.Do(req)
		if err != nil {
			log.Warn().Err(err).Msg("failed to get a response")
		} else {
			log.Debug().Msgf("%d", res.StatusCode)
			err := setLiveness(tp.cache, Liveness{URL: url, Status: res.StatusCode})
			if err != nil {
				log.Error().Err(err).Msg("failed to set task")
			}
		}
	}
}

func setLiveness(cache *bigcache.BigCache, liveness Liveness) error {
	bytes, err := cache.Get(LivenessCacheName)
	if err != nil {
		return cache.Set(LivenessCacheName, serializeLiveness(map[string]Liveness{liveness.URL: liveness}))
	} else {
		var livenesses map[string]Liveness
		err = json.Unmarshal(bytes, &livenesses)
		if err != nil {
			return err
		}
		livenesses[liveness.URL] = liveness
		return cache.Set(LivenessCacheName, serializeLiveness(livenesses))
	}
}

func serializeLiveness(liveness map[string]Liveness) []byte {
	bytes, _ := json.Marshal(&liveness)
	return bytes
}
