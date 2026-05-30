package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/model"
	"github.com/HOangAG2207/GoBeK03/internal/service/health/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(t *testing.T) *mocks.Health

		expectedStatus int
		expectedBody   interface{}
		expectError    bool
	}{
		{
			name: "success",

			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)

				mockSvc.
					On("CheckHealth", mock.Anything).
					Return(&model.Health{
						Message:     "OK",
						ServiceName: "test-service",
						InstanceID:  "instance-1",
					}, nil)

				return mockSvc
			},

			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedBody: &model.Health{
				Message:     "OK",
				ServiceName: "test-service",
				InstanceID:  "instance-1",
			},
		},
		{
			name: "service error",

			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)

				mockSvc.
					On("CheckHealth", mock.Anything).
					Return(nil, errors.New("some error"))

				return mockSvc
			},

			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
			expectedBody: map[string]interface{}{
				"error": "Internal Server Error!",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// ===== 1. Setup Echo =====
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// ===== 2. Mock =====
			mockSvc := tc.setupMock(t)

			// ===== 3. Handler =====
			h := &healthHandler{
				service: mockSvc,
			}

			// ===== 4. Call =====
			err := h.CheckHealth(ctx)

			// ===== 5. Assert error =====
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// ===== 6. Assert status =====
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// ===== 7. Assert body =====
			expectedJSON, _ := json.Marshal(tc.expectedBody)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
