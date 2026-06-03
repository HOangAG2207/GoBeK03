package service

import (
	"context"

	"github.com/HOangAG2207/GoBeK03/internal/model"
	repository "github.com/HOangAG2207/GoBeK03/internal/repository/health"
)

const (
	StatusOK         = "OK"
	RedisPingTimeout = "redis: client is closed"
	DBPingConfused   = "database: database is closed"
)

//go:generate mockery --name Health --filename check_health_mock.go --output ./mocks
type Health interface {

	// CheckHealth xử lý logic kiểm tra trạng thái hệ thống
	// Trả về:
	// - model.Health: thông tin health status
	// - error: lỗi nếu có trong quá trình xử lý
	CheckHealth(ctx context.Context) (*model.Health, error)
}

// healthService là implementation cụ thể của interface Health
// Nó chứa business logic thật sự của health-check
type healthService struct {

	// repo là dependency tầng repository
	// dùng để kiểm tra trạng thái hạ tầng (Redis, DB,...)
	repo repository.Health

	// serviceName là tên service (inject từ config)
	serviceName string

	// instanceId là định danh instance của service
	// có thể được truyền từ config hoặc generate runtime
	instanceId string
}

// NewHealth là constructor của service Health
// Dùng dependency injection để truyền repo + metadata vào service
func NewHealthService(repo repository.Health, serviceName, instanceId string) Health {

	// Trả về interface Health thay vì struct cụ thể
	// giúp ẩn implementation và dễ mock/test hơn
	return &healthService{
		repo:        repo,
		serviceName: serviceName,
		instanceId:  instanceId,
	}
}
