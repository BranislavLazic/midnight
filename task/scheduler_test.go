package task_test

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api/testapi"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/branislavlazic/midnight/task"
	"github.com/go-co-op/gocron"
	"strconv"
	"testing"
	"time"
)

func TestSchedulerRunAll(t *testing.T) {
	serviceRepo := postgres.NewServiceRepository(testapi.DB)
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to initialize cache %s", err.Error())
	}
	err = serviceRepo.DeleteAll()
	if err != nil {
		t.Fatalf("failed to delete all services")
	}
	serviceID := model.ServiceID(1)
	service := &model.Service{ID: serviceID, Name: "test service", URL: "http://testtest1", CheckIntervalSeconds: 5}
	_, err = serviceRepo.Create(service)
	if err != nil {
		t.Fatalf("failed to create the service %s", err.Error())
	}
	scheduler := gocron.NewScheduler(time.UTC)
	taskProvider := task.NewProvider(cache)
	taskScheduler := task.NewScheduler(scheduler, taskProvider, serviceRepo)
	err = taskScheduler.RunAll()
	if err != nil {
		t.Fatalf("failed to run all tasks %s", err.Error())
	}
	if !scheduler.IsRunning() {
		t.Fatalf("scheduler is not running the task")
	}
	jobs, err := scheduler.FindJobsByTag(strconv.FormatInt(int64(serviceID), 10))
	if err != nil {
		t.Fatal("failed to find any jobs")
	}
	if len(jobs) == 0 {
		t.Fatalf("exprected 1 job. got %d", len(jobs))
	}
}

func TestSchedulerAddTask(t *testing.T) {
	serviceRepo := postgres.NewServiceRepository(testapi.DB)
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	if err != nil {
		t.Fatalf("failed to initialize cache %s", err.Error())
	}
	scheduler := gocron.NewScheduler(time.UTC)
	taskProvider := task.NewProvider(cache)
	taskScheduler := task.NewScheduler(scheduler, taskProvider, serviceRepo)
	taskID := int64(1)
	cfg := task.Config{ID: taskID, Name: "test config", URL: "http://testtest1", Timeout: 10}
	err = taskScheduler.Add(cfg, 10)
	if err != nil {
		t.Fatalf("failed to add the task %s", err.Error())
	}
	if !scheduler.IsRunning() {
		t.Fatal("scheduler is not running the task")
	}
	jobs, err := scheduler.FindJobsByTag(strconv.FormatInt(taskID, 10))
	if err != nil {
		t.Fatal("failed to find any jobs")
	}
	if len(jobs) == 0 {
		t.Fatalf("exprected 1 job. got %d", len(jobs))
	}
}
