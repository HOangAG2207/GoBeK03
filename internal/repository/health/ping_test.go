package repository

import (
	"testing"
)

// TestRepo_HealthPing kiểm tra hàm HealthPing của repository Health
// nhằm đảm bảo repository trả về kết quả đúng như kỳ vọng
func TestRepo_HealthPing(t *testing.T) {

	// Danh sách test cases theo kiểu table-driven test
	// giúp dễ mở rộng nhiều scenario sau này
	tests := []struct {

		// name là tên của từng test case
		// dùng để hiển thị khi chạy t.Run
		name string

		// expected là kết quả mong đợi của hàm HealthPing()
		expected bool
	}{
		{
			// Test case 1: trường hợp ping thành công
			name:     "Ping OK",
			expected: true,
		},
	}

	// Khởi tạo repository cần test
	// NewHealth() trả về implementation của Health interface
	repo := NewHealth()

	// Duyệt qua từng test case
	for _, tt := range tests {

		// t.Run giúp tách từng test case riêng biệt
		// thuận tiện debug khi fail
		t.Run(tt.name, func(t *testing.T) {

			// Gọi hàm cần test
			result := repo.HealthPing()

			// So sánh kết quả thực tế với kết quả mong đợi
			if result != tt.expected {

				// Nếu không khớp → báo lỗi test
				t.Errorf(
					"HealthPing() = %v, want %v",
					result,
					tt.expected,
				)
			}
		})
	}
}
