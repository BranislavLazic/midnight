package api_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/branislavlazic/midnight/api/testapi"
)

func TestEnvironment(t *testing.T) {
	scenarios := []testapi.ApiScenario{
		{
			Name:            "empty data",
			Method:          http.MethodGet,
			Url:             "/v1/environments",
			ExpectedStatus:  http.StatusOK,
			Body:            strings.NewReader(``),
			ExpectedContent: []string{`[]`},
		},
		{
			Name:           "create",
			Method:         http.MethodPost,
			Url:            "/v1/environments",
			ExpectedStatus: http.StatusCreated,
			Body:           strings.NewReader(`{"name":"PROD"}`),
			RequestHeaders: testapi.LongLivedAuthorizationHeader,
		},
	}
	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}

// func TestEnvironmentFound(t *testing.T) {
// 	envRepo := postgres.NewEnvironmentRepository(testapi.DB)
// 	_ = envRepo.DeleteAll()
// 	_, err := envRepo.Create(&model.Environment{ID: model.EnvironmentID(1), Name: "PROD"})
// 	if err != nil {
// 		t.Fatal("failed to create an environment")
// 	}
// 	app := testapi.InitTestApp()
// 	req := httptest.NewRequest(fiber.MethodGet, "/v1/environments", nil)
// 	res, err := app.Test(req)
// 	utils.AssertEqual(t, nil, err, "app.Test(req)")
// 	utils.AssertEqual(t, 200, res.StatusCode, "Status code")

// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatalf("failed to read response body")
// 	}
// 	expectedBody := `[{"id":1,"name":"PROD"}]`
// 	utils.AssertEqual(t, expectedBody, string(body), "Body")
// }
