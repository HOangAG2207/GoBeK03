package url

import (
	"context"
	"errors"
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
				return pkgredis.InitMockRedis(t)
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
				return pkgredis.InitMockRedis(t)
			},

			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// ===== 1. Setup redis =====
			redisClient := tc.setupMock(ctx)

			// ===== 2. Init repo =====
			repo := NewUrlRepository(redisClient, tc.expireIn)

			// ===== 3. Call function =====
			result, err := repo.StoreURLIfAbsent(ctx, tc.code, tc.url)

			// ===== 4. Assert result =====
			assert.Equal(t, tc.expectedResult, result)

			// ===== 5. Assert error (FIX QUAN TRỌNG) =====
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
