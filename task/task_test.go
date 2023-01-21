package task_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/task"
	"github.com/rs/zerolog/log"
)

func TestNewTask(t *testing.T) {
	const taskID = 1
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to initialize cache %s", err.Error())
	}
	provider := task.NewProvider(cache)
	cfg := task.Config{ID: taskID, Name: "Test service", URL: "http://testservice", Timeout: 10}
	taskFn := provider.NewTask(cfg)
	// invoke
	taskFn()
	bytes, err := cache.Get(task.ServiceStatusCacheName)
	if err != nil {
		t.Fatal("failed to read the cache")
	}

	var serviceStatuses map[int64]task.ServiceStatus
	_ = json.Unmarshal(bytes, &serviceStatuses)
	serviceStatus := serviceStatuses[taskID]
	if serviceStatus.StatusCode != 404 {
		t.Fatalf("expected status code 404. got %d", serviceStatus.StatusCode)
	}
	if serviceStatus.Name != "Test service" {
		t.Fatalf("expected name 'Test service'. got %s", serviceStatus.Name)
	}
	if serviceStatus.URL != "http://testservice" {
		t.Fatalf("expected name 'http://testservice'. got %s", serviceStatus.URL)
	}
}

func TestRemoveTask(t *testing.T) {
	const taskID = 1
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize cache")
	}
	provider := task.NewProvider(cache)
	cfg := task.Config{ID: taskID, Name: "Test service", URL: "http://testservice", Timeout: 10}
	taskFn := provider.NewTask(cfg)
	// invoke
	taskFn()
	bytes, err := cache.Get(task.ServiceStatusCacheName)
	if err != nil {
		t.Fatal("failed to read the cache")
	}
	if string(bytes) == "{}" {
		t.Fatal("the cache should not be empty")
	}
	// remove
	err = provider.RemoveTask(taskID)
	if err != nil {
		t.Fatal("failed to remove the task")
	}
	bytesAfterRemoval, err := cache.Get(task.ServiceStatusCacheName)
	if err != nil {
		t.Fatal("failed to read the cache")
	}
	if string(bytesAfterRemoval) != "{}" {
		t.Fatal("the cache should be empty")
	}
}
