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
)

func TestHandler_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(t *testing.T) *mocks.Health

		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)
				mockSvc.On("CheckHealth").Return(&model.Health{
					Message:     "OK",
					ServiceName: "test-service",
					InstanceID:  "instance-1",
				}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
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
				mockSvc.On("CheckHealth").Return(nil, errors.New("some error"))
				return mockSvc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Internal Server Error!",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// setup echo
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			// setup mock service
			mockSvc := tc.setupMock(t)

			// init handler
			h := &healthHandler{
				service: mockSvc,
			}

			// call handler
			err := h.CheckHealth(ctx)

			// assert error (Echo handler thường return err khi fail)
			if tc.expectedStatus == http.StatusInternalServerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// assert status
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// assert body
			expectedJSON, _ := json.Marshal(tc.expectedBody)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
