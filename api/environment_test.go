package api_test

import (
	"github.com/branislavlazic/midnight/api/testapi"
	"github.com/branislavlazic/midnight/model"
	"github.com/branislavlazic/midnight/repository/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"net/http/httptest"
	"testing"
)

func TestEnvironmentNoData(t *testing.T) {
	envRepo := postgres.NewEnvironmentRepository(testapi.DB)
	_ = envRepo.DeleteAll()
	app := testapi.InitTestApp()
	req := httptest.NewRequest(fiber.MethodGet, "/v1/environments", nil)
	res, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

	body, err := io.ReadAll(res.Body)
	expectedBody := `[]`
	utils.AssertEqual(t, expectedBody, string(body), "Body")
}

func TestEnvironmentFound(t *testing.T) {
	envRepo := postgres.NewEnvironmentRepository(testapi.DB)
	_ = envRepo.DeleteAll()
	_, err := envRepo.Create(&model.Environment{ID: model.EnvironmentID(1), Name: "PROD"})
	if err != nil {
		t.Fatal("failed to create an environment")
	}
	app := testapi.InitTestApp()
	req := httptest.NewRequest(fiber.MethodGet, "/v1/environments", nil)
	res, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

	body, err := io.ReadAll(res.Body)
	expectedBody := `[{"id":1,"Name":"PROD"}]`
	utils.AssertEqual(t, expectedBody, string(body), "Body")
}
