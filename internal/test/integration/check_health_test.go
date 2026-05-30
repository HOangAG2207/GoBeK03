package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
	"github.com/stretchr/testify/assert"
)

// TestEndpoint_CheckHealth là integration test
// kiểm tra toàn bộ flow: HTTP → handler → service → repository
func TestEndpoint_CheckHealth(t *testing.T) {

	// Cho phép test chạy song song với các test khác
	t.Parallel()

	// ===== 1. Fake config cho test =====

	cfg := &api.Config{
		ServiceName: "test-service",
		InstanceID:  "test-instance",
		AppPort:     "8080",
	}

	// ===== 2. Định nghĩa test cases =====

	testCases := []struct {

		// tên test case
		name string

		// function setup HTTP request + execute engine
		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		// HTTP status mong đợi
		expectedStatusCode int

		// response body mong đợi (dạng JSON string)
		expectedResponseBody string
	}{
		{
			name: "health_check_success",

			// setup request và chạy engine
			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {

				// tạo HTTP request giả lập
				req := httptest.NewRequest(
					http.MethodGet,
					"/api/health-check",
					nil,
				)

				// recorder để capture response
				rec := httptest.NewRecorder()

				// gọi full HTTP server pipeline (không cần start server thật)
				engine.ServeHTTP(rec, req)

				return rec
			},

			// mong đợi HTTP 200 OK
			expectedStatusCode: http.StatusOK,

			// expected JSON response dạng string
			expectedResponseBody: `{
				"message": "OK",
				"service_name": "test-service",
				"instance_id": "test-instance"
			}`,
		},
	}

	// ===== 3. Loop từng test case =====

	for _, tc := range testCases {

		// capture biến tránh lỗi khi chạy parallel
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			// chạy song song sub-test
			t.Parallel()

			// ===== 4. Khởi tạo engine (full application) =====

			engine := api.NewEngine(&api.EngineOpts{
				Cfg:         cfg, // inject config test vào hệ thống
				RedisClient: pkgredis.InitMockRedis(t),
			})

			// ===== 5. Execute HTTP request =====

			rec := tc.setupTestHTTP(engine)

			// ===== 6. Assert HTTP status code =====
			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			// ===== 7. Parse actual response JSON =====
			var actual map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &actual)
			assert.NoError(t, err)

			// ===== 8. Parse expected response JSON =====
			var expected map[string]interface{}
			err = json.Unmarshal([]byte(tc.expectedResponseBody), &expected)
			assert.NoError(t, err)

			// ===== 9. Compare JSON structurally =====

			// so sánh object thay vì string (tránh lỗi format spacing)
			assert.Equal(t, expected, actual)
		})
	}
}
