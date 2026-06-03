package pkgredis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient khởi tạo và trả về Redis client đã được kết nối sẵn
// envprefix: prefix dùng để đọc biến môi trường (ví dụ: REDIS_HOST, REDIS_PORT,...)
func NewRedisClient(envprefix string) (*redis.Client, error) {

	// ===== 1. Load config Redis từ environment =====

	// newConfig sẽ đọc env (theo prefix) và build Redis config
	cfg, err := newConfig(envprefix)

	// nếu không đọc được config → trả error ngay
	if err != nil {
		return nil, err
	}

	// ===== 2. Tạo Redis client =====

	redisClient := redis.NewClient(&redis.Options{

		// địa chỉ Redis server (host:port)
		Addr: cfg.Address,

		// password nếu Redis có bật auth
		Password: cfg.Password,

		// database index (Redis hỗ trợ nhiều DB 0,1,2,...)
		DB: cfg.DB,
	})

	// ===== 3. Kiểm tra kết nối Redis =====

	// Ping Redis để verify connection ngay khi startup
	// giúp fail-fast nếu Redis không reachable
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	// ===== 4. Trả về Redis client =====

	return redisClient, nil
}
