package api_test

import (
	"github.com/branislavlazic/midnight/api/testapi"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"net/http/httptest"
	"testing"
)

func TestServiceStatus(t *testing.T) {
	app := testapi.InitTestApp()
	res, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/v1/status", nil))
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, res.StatusCode, "Status code")
}
