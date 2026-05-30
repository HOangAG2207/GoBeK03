package service

import (
	"context"
	"errors"
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/repository/health/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CheckHealth(t *testing.T) {
	t.Parallel() // Cho phép test chạy song song

	// ===== Danh sách test case =====
	testCases := []struct {
		name string

		serviceName string // tên service (inject vào service)
		instanceID  string // instance ID ban đầu

		setupRepoMock func(m *mocks.Health) // setup behavior cho repo mock

		expectNil      bool   // có mong đợi result = nil không
		expectInstance string // giá trị instance mong đợi ("not-empty" = chỉ cần khác rỗng)
		expectError    bool   // có mong đợi error không
	}{
		{
			name:        "repo error -> return nil",
			serviceName: "test-service",
			instanceID:  "instance-1",

			// ===== Case: Redis lỗi =====
			setupRepoMock: func(m *mocks.Health) {
				// Khi gọi RedisPing → trả về lỗi
				m.On("RedisPing", mock.Anything).
					Return(errors.New("redis down"))
			},

			expectNil:   true, // result phải nil
			expectError: true, // phải có error
		},
		{
			name:        "repo healthy with predefined instanceID",
			serviceName: "test-service",
			instanceID:  "instance-1",

			// ===== Case: Redis OK =====
			setupRepoMock: func(m *mocks.Health) {
				m.On("RedisPing", mock.Anything).
					Return(nil)
			},

			expectNil:      false,
			expectInstance: "instance-1", // giữ nguyên instanceID
		},
		{
			name:        "repo healthy with empty instanceID -> generate new",
			serviceName: "test-service",
			instanceID:  "",

			// ===== Case: Redis OK nhưng instanceID rỗng =====
			setupRepoMock: func(m *mocks.Health) {
				m.On("RedisPing", mock.Anything).
					Return(nil)
			},

			expectNil:      false,
			expectInstance: "not-empty", // phải generate ID mới
		},
	}

	// ===== Lặp qua từng test case =====
	for _, tc := range testCases {
		tc := tc // capture biến để tránh lỗi khi chạy parallel

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // mỗi test case chạy song song

			// ===== 1. Setup mock repository =====
			repoMock := mocks.NewHealth(t)

			// Apply behavior mock theo từng case
			tc.setupRepoMock(repoMock)

			// ===== 2. Khởi tạo service =====
			s := &healthService{
				repo:        repoMock,
				serviceName: tc.serviceName,
				instanceId:  tc.instanceID,
			}

			// ===== 3. Gọi hàm cần test =====
			ctx := context.Background()
			result, err := s.CheckHealth(ctx)

			// ===== 4. Assert error =====
			if tc.expectError {
				assert.Error(t, err) // phải có lỗi
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)

			// ===== 5. Assert nil result =====
			if tc.expectNil {
				assert.Nil(t, result)
				return
			}

			// ===== 6. Assert dữ liệu trả về =====
			assert.NotNil(t, result)

			// Message luôn phải là "OK" khi thành công
			assert.Equal(t, "OK", result.Message)

			// ServiceName phải đúng với config
			assert.Equal(t, tc.serviceName, result.ServiceName)

			// ===== 7. Assert instance ID =====
			if tc.expectInstance == "not-empty" {
				// Trường hợp generate ID mới
				assert.NotEmpty(t, result.InstanceID)
			} else {
				// Trường hợp giữ nguyên ID
				assert.Equal(t, tc.expectInstance, result.InstanceID)
			}
		})
	}
}
