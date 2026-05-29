package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestEndpoint_CheckHealth(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		ServiceName: "test-service",
		InstanceID:  "test-instance",
		AppPort:     "8080",
	}

	testCases := []struct {
		name                 string
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "health_check_success",
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/api/health-check", nil)
				rec := httptest.NewRecorder()

				engine.ServeHTTP(rec, req)

				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"message": "OK",
				"service_name": "test-service",
				"instance_id": "test-instance"
			}`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine := api.NewEngine(&api.EngineOpts{
				Cfg: cfg,
			})

			rec := tc.setupTestHTTP(engine)

			// ✅ Status code
			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			// ✅ Convert actual response -> map
			var actual map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &actual)
			assert.NoError(t, err)

			// ✅ Convert expected string -> map
			var expected map[string]interface{}
			err = json.Unmarshal([]byte(tc.expectedResponseBody), &expected)
			assert.NoError(t, err)

			// ✅ Compare JSON logically
			assert.Equal(t, expected, actual)
		})
	}
}
