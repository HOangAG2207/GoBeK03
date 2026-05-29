package service

import (
	"testing"

	"github.com/HOangAG2207/GoBeK03/internal/repository/health/mocks"
	"github.com/stretchr/testify/assert"
)

// TestService_CheckHealth kiểm tra business logic của health service
// Service layer sẽ dựa vào repository để quyết định trả về health status
func TestService_CheckHealth(t *testing.T) {

	// Cho phép test chạy song song với các test khác
	t.Parallel()

	// Table-driven test: định nghĩa nhiều scenario khác nhau
	testCases := []struct {

		// name: tên test case để dễ debug khi fail
		name string

		// serviceName: tên service truyền vào service layer
		serviceName string

		// instanceID: id của instance (có thể được generate hoặc truyền vào)
		instanceID string

		// setupRepoMock: hàm setup behavior cho repository mock
		setupRepoMock func(m *mocks.Health)

		// expectNil: xác định result có phải nil hay không
		expectNil bool

		// expectInstance: giá trị instanceID mong đợi
		// có thể là:
		// - giá trị cố định
		// - "not-empty" để chỉ check khác rỗng
		expectInstance string
	}{
		{
			name: "repo not healthy -> return nil",

			// input service name
			serviceName: "test-service",

			// input instance id
			instanceID: "instance-1",

			// mock repository trả về false (system unhealthy)
			setupRepoMock: func(m *mocks.Health) {
				m.On("HealthPing").Return(false)
			},

			// expected: service trả về nil (không tạo health response)
			expectNil: true,
		},
		{
			name:        "repo healthy with predefined instanceID",
			serviceName: "test-service",
			instanceID:  "instance-1",

			// mock repo trả về healthy
			setupRepoMock: func(m *mocks.Health) {
				m.On("HealthPing").Return(true)
			},

			// expect service trả về object không nil
			expectNil: false,

			// instance phải giữ nguyên value truyền vào
			expectInstance: "instance-1",
		},
		{
			name:        "repo healthy with empty instanceID -> generate uuid",
			serviceName: "test-service",
			instanceID:  "",

			// repo vẫn healthy
			setupRepoMock: func(m *mocks.Health) {
				m.On("HealthPing").Return(true)
			},

			// service phải generate instanceID mới
			expectNil: false,

			// kiểm tra instanceID được generate (không được rỗng)
			expectInstance: "not-empty",
		},
	}

	// Lặp qua từng test case
	for _, tc := range testCases {

		// capture biến để tránh bug khi chạy parallel
		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			// cho phép chạy song song từng sub-test
			t.Parallel()

			// ===== 1. Setup repository mock =====

			repoMock := mocks.NewHealth(t)

			// gán behavior cho mock repo theo từng test case
			tc.setupRepoMock(repoMock)

			// ===== 2. Khởi tạo service =====

			s := &healthService{

				// inject mock repository vào service
				repo: repoMock,

				// inject metadata
				serviceName: tc.serviceName,
				instanceId:  tc.instanceID,
			}

			// ===== 3. Gọi function cần test =====

			result, err := s.CheckHealth()

			// service không được trả error trong test case này
			assert.NoError(t, err)

			// ===== 4. Assert nil case =====

			if tc.expectNil {

				// nếu expect nil → result phải nil
				assert.Nil(t, result)
				return
			}

			// ===== 5. Assert result không nil =====

			assert.NotNil(t, result)

			// message luôn phải là OK khi healthy
			assert.Equal(t, "OK", result.Message)

			// service name phải khớp input
			assert.Equal(t, tc.serviceName, result.ServiceName)

			// ===== 6. Assert instance ID =====

			if tc.expectInstance == "not-empty" {

				// case generate uuid → chỉ check không rỗng
				assert.NotEmpty(t, result.InstanceID)
			} else {

				// case cố định → phải đúng giá trị truyền vào
				assert.Equal(t, tc.expectInstance, result.InstanceID)
			}
		})
	}
}
