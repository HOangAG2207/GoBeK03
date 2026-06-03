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
	t.Parallel() // Cho phép test này chạy song song với các test khác

	// ===== Danh sách các test case =====
	testCases := []struct {
		name string

		// Hàm setup mock service tương ứng với từng test case
		setupMock func(t *testing.T) *mocks.Health

		expectedStatus int         // HTTP status mong đợi
		expectedBody   interface{} // Body response mong đợi
		expectError    bool        // Có mong đợi handler trả về error hay không
	}{
		{
			name: "success",

			// ===== Mock cho case thành công =====
			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)

				// Khi gọi CheckHealth với bất kỳ context nào
				// → trả về dữ liệu health hợp lệ và không có lỗi
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

			// ===== Mock cho case service trả lỗi =====
			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)

				// Khi gọi CheckHealth → trả về lỗi
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

	// ===== Lặp qua từng test case =====
	for _, tc := range testCases {
		tc := tc // capture biến để tránh lỗi khi chạy parallel

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // mỗi test case chạy song song

			// ===== 1. Setup Echo context =====
			e := echo.New()

			// Tạo request GET /health
			req := httptest.NewRequest(http.MethodGet, "/health", nil)

			// Recorder để ghi lại response
			rec := httptest.NewRecorder()

			// Tạo context từ Echo
			ctx := e.NewContext(req, rec)

			// ===== 2. Setup mock service =====
			mockSvc := tc.setupMock(t)

			// ===== 3. Khởi tạo handler với mock service =====
			h := &healthHandler{
				service: mockSvc,
			}

			// ===== 4. Gọi handler =====
			err := h.CheckHealth(ctx)

			// ===== 5. Assert error =====
			if tc.expectError {
				assert.Error(t, err) // mong đợi có lỗi
			} else {
				assert.NoError(t, err) // không mong đợi lỗi
			}

			// ===== 6. Assert HTTP status code =====
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// ===== 7. Assert response body =====
			// Convert expected body sang JSON để so sánh
			expectedJSON, _ := json.Marshal(tc.expectedBody)

			// So sánh JSON (không phụ thuộc thứ tự field)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
