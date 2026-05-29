package repository

//go:generate mockery --name Health --filename check_health_mock.go --output ./mocks
type Health interface {

	// HealthPing kiểm tra trạng thái của hệ thống (DB, Redis, service...)
	// Trả về true nếu OK, false nếu có vấn đề
	HealthPing() bool
}

// healthRepository là implementation cụ thể của interface Health
type healthRepository struct {
	// Hiện tại chưa có dependency (DB, Redis,...)
	// Sau này có thể inject thêm vào đây
}

// NewHealth là constructor tạo mới repository
func NewHealth() Health {

	// Trả về dưới dạng interface → giúp loose coupling
	return &healthRepository{}
}
