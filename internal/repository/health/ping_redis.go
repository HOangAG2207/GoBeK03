package repository

import "context"

// HealthPing kiểm tra trạng thái "sức khỏe" của hệ thống ở tầng repository
// Kiểm tra ping kết nối đến redis
func (h *healthRepository) RedisPing(ctx context.Context) error {
	return h.redis.Ping(ctx).Err()
}
