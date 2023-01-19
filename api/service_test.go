package api_test

import (
	"github.com/branislavlazic/midnight/api/session"
	"github.com/branislavlazic/midnight/api/testapi"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceNoData(t *testing.T) {
	serviceRepo := postgres.NewServiceRepository(testapi.DB)
	_ = serviceRepo.DeleteAll()
	app := testapi.InitTestApp()
	sessionCookieID := testapi.GenerateSecureSession(t)
	req := httptest.NewRequest(fiber.MethodGet, "/v1/services", nil)
	req.AddCookie(&http.Cookie{Name: session.SecureCookieName, Value: sessionCookieID})
	res, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

	body, err := io.ReadAll(res.Body)
	expectedBody := `[]`
	utils.AssertEqual(t, expectedBody, string(body), "Body")
}

func TestServiceFound(t *testing.T) {
	const serviceID = model.ServiceID(1)
	serviceRepo := postgres.NewServiceRepository(testapi.DB)
	_ = serviceRepo.DeleteAll()
	_, err := serviceRepo.Create(&model.Service{ID: serviceID, Name: "Test service", URL: "http://service.com", CheckIntervalSeconds: 30})
	if err != nil {
		t.Fatal("failed to create a service")
	}
	app := testapi.InitTestApp()
	sessionCookieID := testapi.GenerateSecureSession(t)
	req := httptest.NewRequest(fiber.MethodGet, "/v1/services", nil)
	req.AddCookie(&http.Cookie{Name: session.SecureCookieName, Value: sessionCookieID})
	res, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

	body, err := io.ReadAll(res.Body)
	expectedBody := `[{"id":1,"name":"Test service","url":"http://service.com","checkIntervalSeconds":30}]`
	utils.AssertEqual(t, expectedBody, string(body), "Body")
}
