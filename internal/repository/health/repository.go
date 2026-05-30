package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

//go:generate mockery --name Health --filename check_health_mock.go --output ./mocks
type Health interface {

	// HealthPing kiểm tra trạng thái của hệ thống (DB, Redis, service...)
	// Trả về true nếu OK, false nếu có vấn đề
	RedisPing(ctx context.Context) error
}

// healthRepository là implementation cụ thể của interface Health
type healthRepository struct {
	// Hiện tại chưa có dependency (DB, Redis,...)
	// Sau này có thể inject thêm vào đây
	redis *redis.Client
}

// NewHealth là constructor tạo mới repository
func NewHealthRepository(redisClient *redis.Client) Health {

	// Trả về dưới dạng interface → giúp loose coupling
	return &healthRepository{
		redis: redisClient,
	}
}
