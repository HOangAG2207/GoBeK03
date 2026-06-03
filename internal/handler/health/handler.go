package handler

import (
	service "github.com/HOangAG2207/GoBeK03/internal/service/health"
	"github.com/labstack/echo/v4"
)

// Health là interface định nghĩa các handler liên quan đến health-check
// Giúp:
// - Dễ mock khi test
// - Không phụ thuộc vào implementation cụ thể
type Health interface {

	// CheckHealth xử lý HTTP request (Echo context)
	// Trả về error để Echo xử lý nếu có lỗi
	CheckHealth(ctx echo.Context) error
}

// healthHandler là struct implement interface Health
type healthHandler struct {

	// service là dependency của handler (business logic layer)
	// Interface này nằm ở layer service → đúng clean architecture
	service service.Health
}

// NewHealth là constructor để khởi tạo handler
// Nhận vào service (dependency injection)
func NewHealthHandler(svc service.Health) Health {

	// Trả về dưới dạng interface (Health)
	// → giúp ẩn implementation bên trong
	return &healthHandler{
		service: svc,
	}
}
