package api_test

import (
	"github.com/branislavlazic/midnight/api/testapi"
	"github.com/branislavlazic/midnight/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"net/http/httptest"
	"testing"
)

func TestServiceStatusNoStatuses(t *testing.T) {
	app := testapi.InitTestApp()
	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/status", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

	body, err := io.ReadAll(res.Body)
	expectedBody := `[]`
	utils.AssertEqual(t, expectedBody, string(body), "Body")
}

func TestServiceStatusSingleStatus(t *testing.T) {
	err := task.SaveServiceStatus(testapi.Cache, task.ServiceStatus{
		ID:         1,
		Name:       "Test service",
		URL:        "http://testservice.com",
		Status:     "200 OK",
		StatusCode: 200,
	})
	if err != nil {
		t.Fatal("failed to save the service status")
	}
	app := testapi.InitTestApp()
	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/status", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

	body, err := io.ReadAll(res.Body)
	expectedBody := `[{"id":1,"name":"Test service","url":"http://testservice.com","status":"200 OK","statusCode":200}]`
	utils.AssertEqual(t, expectedBody, string(body), "Body")
}
