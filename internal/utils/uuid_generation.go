package utils

import "github.com/google/uuid"

// UuidGenerator tạo một UUID mới theo chuẩn RFC 4122
// và trả về dưới dạng string
func UuidGenerator() string {

	// uuid.New() tạo một UUID version 4 (random-based UUID)
	// đảm bảo gần như không trùng lặp trong thực tế
	return uuid.New().String()
}
