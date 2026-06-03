package repository

import (
	"testing"

	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_RedisPing(t *testing.T) {
	t.Parallel() // Cho phép toàn bộ test này chạy song song

	// ===== Danh sách test case =====
	testCases := []struct {
		name string

		// Hàm setup mock redis client cho từng case
		setupMock func() *redis.Client

		expectedError error // Lỗi mong đợi khi gọi RedisPing
	}{
		{
			name: "successful ping",

			// ===== Case: Redis hoạt động bình thường =====
			setupMock: func() *redis.Client {
				// Khởi tạo mock Redis (thường là miniredis hoặc tương tự)
				redisClient := pkgredis.InitMockRedis(t)
				return redisClient
			},

			// Không có lỗi khi ping thành công
			expectedError: nil,
		},
		{
			name: "failed ping",

			// ===== Case: Redis bị đóng kết nối =====
			setupMock: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)

				// Đóng client → giả lập Redis bị down / mất kết nối
				redisClient.Close()

				return redisClient
			},

			// Mong đợi lỗi ErrClosed từ redis client
			expectedError: redis.ErrClosed,
		},
	}

	// ===== Lặp qua từng test case =====
	for _, tc := range testCases {
		tc := tc // capture biến để tránh lỗi khi chạy parallel

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // mỗi test case chạy song song

			// Context dùng cho Redis (Go 1.21+ có t.Context())
			ctx := t.Context()

			// ===== 1. Setup Redis mock =====
			redisMockClient := tc.setupMock()

			// ===== 2. Khởi tạo repository =====
			healthCheckRepo := NewHealthRepository(redisMockClient)

			// ===== 3. Gọi hàm cần test =====
			err := healthCheckRepo.RedisPing(ctx)

			// ===== 4. Assert kết quả =====
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
