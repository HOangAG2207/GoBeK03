package repository

// HealthPing kiểm tra trạng thái "sức khỏe" của hệ thống ở tầng repository
func (repo *healthRepository) HealthPing() bool {

	// Hiện tại luôn trả về true
	// → nghĩa là hệ thống luôn được coi là "healthy"
	return true
}
