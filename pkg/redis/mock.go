package pkgredis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
}
func InitClosedRedis(t *testing.T) *redis.Client {
	t.Helper()
	return redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:1",
	})
}
