package model

// Health là model dùng để biểu diễn trạng thái "health" của service
// Thường được trả về ở endpoint /health-check
type Health struct {

	// Message: thông điệp trạng thái (ví dụ: "OK", "Healthy", "Down")
	Message string `json:"message"`

	// ServiceName: tên của service (giúp nhận diện khi có nhiều service)
	ServiceName string `json:"service_name"`

	// InstanceID: định danh instance cụ thể của service
	// Hữu ích khi chạy nhiều instance (load balancing, Kubernetes, etc.)
	InstanceID string `json:"instance_id"`
}
