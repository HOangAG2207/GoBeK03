package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	"github.com/HOangAG2207/GoBeK03/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestEndpoint_CheckHealth(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		ServiceName: "test-service",
		InstanceID:  "test-instance",
		AppPort:     "8080",
	}

	testCases := []struct {
		name string

		setupTestHTTP func() *httptest.ResponseRecorder

		expectedStatusCode int
	}{
		{
			name: "health_check_success",
			setupTestHTTP: func() *httptest.ResponseRecorder {
				engine := api.NewEngine(&api.EngineOpts{
					Cfg: cfg,
				})

				req := httptest.NewRequest(http.MethodGet, "/api/health-check", nil)
				rec := httptest.NewRecorder()

				engine.ServeHTTP(rec, req)

				return rec
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.setupTestHTTP()

			// ✅ Check status code
			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			// ✅ Parse JSON response (best practice)
			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			// ✅ Assert đúng theo json tag
			assert.Equal(t, "OK", resp["message"])
			assert.Equal(t, "test-service", resp["service_name"])
			assert.Equal(t, "test-instance", resp["instance_id"])
		})
	}
}
