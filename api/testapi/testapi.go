package testapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"github.com/branislavlazic/midnight/api"
	"github.com/branislavlazic/midnight/api/session"
	"github.com/branislavlazic/midnight/cache"
	"github.com/branislavlazic/midnight/config"
	"github.com/branislavlazic/midnight/db"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
	"time"
)

var DB *gorm.DB
var Cache cache.Internal

const email = "john.doe@mail.com"
const password = "john"
const role = "ADMIN"

type Session struct {
	K string `gorm:"type:VARCHAR(255);primaryKey"`
	V []byte `gorm:"type:BYTEA"`
	E int64  `gorm:"type:BIGINT"`
}

func init() {
	cfg := TestConfig()
	DB, _ = db.GetPGDBPool(cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)
	_ = DB.AutoMigrate(&model.User{})
	_ = DB.AutoMigrate(&model.Service{})
	_ = DB.AutoMigrate(&Session{})
	_ = DB.AutoMigrate(&model.Environment{})
	Cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(24*time.Hour))
}

func TestConfig() *config.AppConfig {
	return &config.AppConfig{
		DbHost:        "localhost",
		DbPort:        5433,
		DbUser:        "postgres",
		DbPassword:    "postgres",
		DbName:        "midnight",
		EnableSwagger: false,
	}
}

func InitTestApp() *fiber.App {
	serverSettings := api.ServerSettings{Config: TestConfig(), DB: DB, Cache: Cache}
	return api.InitApp(serverSettings)
}

func GenerateSecureSession(t *testing.T) string {
	userRepo := postgres.NewUserRepository(DB)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		t.Fatal("failed to create a password")
	}
	_, err = userRepo.Create(&model.User{Email: email, Password: string(hash), Role: role, Enabled: true})
	if err != nil {
		t.Fatal("failed to create the user")
	}
	app := InitTestApp()
	loginReq := model.LoginRequest{Email: email, Password: password}
	b, _ := json.Marshal(&loginReq)
	jsonBody := bytes.NewReader(b)
	_ = json.Unmarshal(b, &model.LoginRequest{})

	req := httptest.NewRequest(fiber.MethodPost, "/v1/login", jsonBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal("failed to login user")
	}
	var securedCookieValue string
	for _, c := range resp.Cookies() {
		if c.Name == session.SecureCookieName {
			securedCookieValue = c.Value
		}
	}
	return securedCookieValue
}
