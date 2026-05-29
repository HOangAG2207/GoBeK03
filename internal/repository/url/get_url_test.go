package url

import (
	"context"
	"testing"

	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client

		expectedError error

		verifyFunc func(ctx context.Context, redisClient *redis.Client)
	}{
		{
			name: "successful URL get",

			setupMock: func() *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				redisClient.Set(t.Context(), "test", "https://example.com", 10000)
				return redisClient
			},

			expectedError: nil,

			verifyFunc: func(ctx context.Context, redisClient *redis.Client) {
				url, err := redisClient.Get(ctx, "test").Result()

				assert.Nil(t, err)

				assert.Equal(t, "https://example.com", url)
			},
		},
		{
			name: "failed URL storage due to closed Redis client",

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

			urlStorage := NewUrlRepository(redisMockClient, 0)

			_, err := urlStorage.GetURL(ctx, "test")
			assert.Equal(t, tc.expectedError, err)

			if err == nil {
				tc.verifyFunc(ctx, redisMockClient)
			}

		})
	}
}
