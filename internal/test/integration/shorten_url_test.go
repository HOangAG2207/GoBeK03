package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestEndpoint_ShortenUrl(t *testing.T) {
	t.Parallel()

	type requestBody struct {
		URL string `json:"url"`
		Exp int64  `json:"exp,omitempty"`
	}

	testCases := []struct {
		name string

		setupMockRedis func(ctx context.Context, redisClient *redis.Client) *redis.Client
		setupTestHTTP  func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode int
		expectCode         bool
	}{
		{
			name: "success - with exp",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return redisClient
			},

			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := requestBody{
					URL: "https://google.com",
					Exp: 3600,
				}
				jsonBody, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode: http.StatusOK,
			expectCode:         true,
		},
		{
			name: "success - default exp",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return redisClient
			},

			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := requestBody{
					URL: "https://google.com",
					// không truyền exp
				}
				jsonBody, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode: http.StatusOK,
			expectCode:         true,
		},
		{
			name: "bad request - empty url",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return redisClient
			},

			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := requestBody{
					URL: "",
				}
				jsonBody, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode: http.StatusBadRequest,
			expectCode:         false,
		},
		{
			name: "bad request - invalid exp",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return redisClient
			},

			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := requestBody{
					URL: "https://google.com",
					Exp: -1,
				}
				jsonBody, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode: http.StatusBadRequest,
			expectCode:         false,
		},
		{
			name: "internal error - redis down",

			setupMockRedis: func(ctx context.Context, redisClient *redis.Client) *redis.Client {
				return pkgredis.InitClosedRedis(t)
			},

			setupTestHTTP: func(engine api.Engine) *httptest.ResponseRecorder {
				body := requestBody{
					URL: "https://google.com",
					Exp: 3600,
				}
				jsonBody, _ := json.Marshal(body)

				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				engine.ServeHTTP(rec, req)

				return rec
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectCode:         false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// ===== 1. Setup redis =====
			mockRedis := pkgredis.InitMockRedis(t)
			redisClient := tc.setupMockRedis(ctx, mockRedis)

			// ===== 2. Setup engine =====
			engine := api.NewEngine(&api.EngineOpts{
				Cfg: &api.Config{
					AppPort:     "8080",
					ServiceName: "test-service",
					InstanceID:  "test-instance",
				},
				RedisClient: redisClient,
			})

			// ===== 3. Call HTTP =====
			rec := tc.setupTestHTTP(engine)

			// ===== 4. Assert status =====
			assert.Equal(t, tc.expectedStatusCode, rec.Code)

			// ===== 5. Assert response =====
			var resp map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			code, _ := resp["code"].(string)

			if tc.expectCode {
				assert.NotEmpty(t, code)

				assert.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]+$`), code)

				assert.GreaterOrEqual(t, len(code), 6)
				assert.LessOrEqual(t, len(code), 12)
			} else {
				assert.Empty(t, code)
			}
		})
	}
}
