package testapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api"
	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/config"
	"github.com/branislavlazic/midnight/db"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Cache cache.Internal

var LongLivedAuthorizationHeader = map[string]string{
	"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNhYTMwMmU1LWRlNGItNGY0Zi1iMjQ0LTM5YjE0MDMzZTgxYSIsImVtYWlsIjoiYWRtaW5AYWRtaW4uY29tIiwicm9sZSI6IlJPTEVfQURNSU4iLCJleHAiOjIzMzM5NjUzODR9.dp3536IfMiALac72g3WuoXSbGLgiLB7Txz6E5uCL9iI",
}

func init() {
	cfg := TestConfig()
	DB, _ = db.GetPGDBPool(cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	_ = DB.AutoMigrate(&model.User{})
	_ = DB.AutoMigrate(&model.Service{})
	Cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
	_ = DB.AutoMigrate(&model.Environment{})
}

func TestConfig() *config.AppConfig {
	return &config.AppConfig{
		DbHost:        "localhost",
		DbPort:        5433,
		DbUser:        "postgres",
		DbPassword:    "postgres",
		DbName:        "midnight",
		EnableSwagger: false,
		SessionSecret: "secret",
	}
}

type ApiScenario struct {
	Name               string
	Method             string
	Url                string
	Body               io.Reader
	RequestHeaders     map[string]string
	ExpectedStatus     int
	ExpectedContent    []string
	NotExpectedContent []string
}

func (scenario *ApiScenario) Test(t *testing.T) {
	settings := api.ServerSettings{Config: TestConfig(), DB: DB, Cache: Cache}
	a := api.NewApp(postgres.NewRepository(DB), settings)
	app := a.InitApi()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(scenario.Method, scenario.Url, scenario.Body)

	app.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancelFunc := context.WithTimeout(c.Request().Context(), 100*time.Millisecond)
			defer cancelFunc()
			c.SetRequest(c.Request().Clone(ctx))
			return next(c)
		}
	})

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range scenario.RequestHeaders {
		req.Header.Set(k, v)
	}
	// execute request
	app.ServeHTTP(recorder, req)
	res := recorder.Result()
	var prefix = scenario.Name
	if prefix == "" {
		prefix = fmt.Sprintf("%s:%s", scenario.Method, scenario.Url)
	}
	if res.StatusCode != scenario.ExpectedStatus {
		t.Errorf("[%s] Expected status code %d, got %d", prefix, scenario.ExpectedStatus, res.StatusCode)
	}
	if len(scenario.ExpectedContent) == 0 && len(scenario.NotExpectedContent) == 0 {
		if len(recorder.Body.Bytes()) != 0 {
			t.Errorf("[%s] Expected empty body, got \n%v", prefix, recorder.Body.String())
		}
	} else {
		buffer := new(bytes.Buffer)
		err := json.Compact(buffer, recorder.Body.Bytes())
		var normalizedBody string
		if err != nil {
			normalizedBody = recorder.Body.String()
		} else {
			normalizedBody = buffer.String()
		}

		for _, item := range scenario.ExpectedContent {
			if !strings.Contains(normalizedBody, item) {
				t.Errorf("[%s] Cannot find %v in response body \n%v", prefix, item, normalizedBody)
				break
			}
		}

		for _, item := range scenario.NotExpectedContent {
			if strings.Contains(normalizedBody, item) {
				t.Errorf("[%s] Didn't expect %v in response body \n%v", prefix, item, normalizedBody)
				break
			}
		}
	}
}
