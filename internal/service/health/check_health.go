package service

import (
	"context"

	"github.com/HOangAG2207/GoBeK03/internal/model"
	"github.com/HOangAG2207/GoBeK03/internal/utils"
)

// CheckHealth là business logic xử lý health-check của hệ thống
// Service sẽ gọi repository để kiểm tra trạng thái hạ tầng
func (s *healthService) CheckHealth(ctx context.Context) (*model.Health, error) {

	// ===== 1. Kiểm tra trạng thái hệ thống qua repository =====

	// repo.HealthPing() thường kiểm tra các dependency như Redis, DB,...
	if err := s.repo.RedisPing(ctx); err != nil {
		return nil, err
	}

	// ===== 2. Xử lý instance ID =====

	// Lấy instanceId từ config (nếu có)
	id := s.instanceId

	// Nếu instanceId rỗng → tự generate UUID
	// đảm bảo mỗi instance luôn có ID duy nhất
	if id == "" {
		id = utils.UuidGenerator()
	}

	// ===== 3. Build response model =====

	// Trả về struct Health chứa thông tin service
	return &model.Health{

		// Message cố định khi hệ thống healthy
		Message: "OK",

		// Tên service lấy từ config
		ServiceName: s.serviceName,

		// Instance ID (config hoặc generated UUID)
		InstanceID: id,
	}, nil
}
