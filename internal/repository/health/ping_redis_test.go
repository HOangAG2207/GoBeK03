package repository

import (
	"testing"

	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_RedisPing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectedError error
	}{
		{
			name: "successful ping",

			setupMock: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				return redisClient
			},

			expectedError: nil,
		},
		{
			name: "failed ping",

			setupMock: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			redisMockClient := tc.setupMock()

			healthCheckRepo := NewHealthRepository(redisMockClient)

			err := healthCheckRepo.RedisPing(ctx)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
