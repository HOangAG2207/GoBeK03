package url

import (
	"context"
	"testing"
	"time"

	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRepository_StoreURLIfAbsent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		code     string
		url      string
		expireIn time.Duration

		setupMock func(ctx context.Context) *redis.Client

		expectedResult bool
		expectedError  error
	}{
		{
			name: "store new URL successfully",

			code:     "abc123",
			url:      "https://example.com",
			expireIn: time.Hour,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				return redisClient
			},

			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "do not overwrite existing URL",

			code:     "abc123",
			url:      "https://example.com/updated",
			expireIn: time.Hour,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				redisClient.Set(ctx, "abc123", "https://example.com", time.Hour)
				return redisClient
			},

			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "failed to store URL due to closed Redis client",

			code:     "abc123",
			url:      "https://example.com",
			expireIn: time.Hour,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				redisClient.Close()
				return redisClient
			},

			expectedResult: false,
			expectedError:  redis.ErrClosed,
		},
		{
			name: "invalid expiration time - negative value",

			code:     "abc123",
			url:      "https://example.com",
			expireIn: -1,

			setupMock: func(ctx context.Context) *redis.Client {
				redisClient := pkgredis.InitMockRedis(t)
				return redisClient
			},

			expectedResult: true, // should still store with default expiration
			expectedError:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			redisClient := tc.setupMock(ctx)
			repo := NewUrlRepository(redisClient, time.Duration(tc.expireIn))

			result, err := repo.StoreURLIfAbsent(ctx, tc.code, tc.url)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
