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
	t.Parallel() // Cho phép test này chạy song song với các test khác

	// Định nghĩa danh sách test case (table-driven test)
	testCases := []struct {
		name string // tên test case

		setupMock func(t *testing.T) *mocks.Health // hàm setup mock service

		expectedStatus int         // HTTP status mong đợi
		expectedBody   interface{} // response body mong đợi
	}{
		{
			name: "success",

			// Setup mock cho trường hợp thành công
			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)

				// Khi gọi CheckHealth → trả về dữ liệu hợp lệ + không lỗi
				mockSvc.On("CheckHealth").Return(&model.Health{
					Message:     "OK",
					ServiceName: "test-service",
					InstanceID:  "instance-1",
				}, nil)

				return mockSvc
			},

			expectedStatus: http.StatusOK,

			// Response mong đợi
			expectedBody: &model.Health{
				Message:     "OK",
				ServiceName: "test-service",
				InstanceID:  "instance-1",
			},
		},
		{
			name: "service error",

			// Setup mock cho trường hợp service trả lỗi
			setupMock: func(t *testing.T) *mocks.Health {
				mockSvc := mocks.NewHealth(t)

				// Khi gọi CheckHealth → trả về lỗi
				mockSvc.On("CheckHealth").Return(nil, errors.New("some error"))

				return mockSvc
			},

			expectedStatus: http.StatusInternalServerError,

			// Response mong đợi khi lỗi
			expectedBody: map[string]interface{}{
				"error": "Internal Server Error!",
			},
		},
	}

	// Lặp qua từng test case
	for _, tc := range testCases {

		tc := tc // capture biến để tránh bug khi chạy parallel

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // mỗi test case chạy song song

			// ===== 1. Setup Echo context =====

			e := echo.New()

			// Tạo HTTP request giả lập
			req := httptest.NewRequest(http.MethodGet, "/health", nil)

			// Recorder để ghi lại response
			rec := httptest.NewRecorder()

			// Tạo context cho Echo
			ctx := e.NewContext(req, rec)

			// ===== 2. Setup mock service =====

			mockSvc := tc.setupMock(t)

			// ===== 3. Khởi tạo handler =====

			h := &healthHandler{
				service: mockSvc, // inject mock vào handler
			}

			// ===== 4. Gọi handler =====

			err := h.CheckHealth(ctx)

			// ===== 5. Assert error =====

			// Nếu expected là lỗi (500) → handler nên return error
			if tc.expectedStatus == http.StatusInternalServerError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// ===== 6. Assert HTTP status =====
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// ===== 7. Assert response body =====

			// Convert expected body → JSON
			expectedJSON, _ := json.Marshal(tc.expectedBody)

			// So sánh JSON (không quan tâm thứ tự field)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
