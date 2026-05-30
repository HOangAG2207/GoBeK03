package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

// charset định nghĩa tập ký tự hợp lệ dùng để generate code
// bao gồm chữ thường, chữ hoa và số
const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateCode tạo một chuỗi random an toàn bằng crypto/rand
// length: độ dài chuỗi cần tạo
// return:
// - string: mã random được sinh ra
// - error: lỗi nếu quá trình random thất bại
func GenerateCode(length int) (string, error) {

	// bytes.Buffer giúp build string hiệu quả hơn + tránh tạo nhiều string tạm
	var strBuilder bytes.Buffer

	// lặp theo độ dài mong muốn
	for i := 0; i < length; i++ {

		// tạo số random trong khoảng [0, len(charset))
		// crypto/rand đảm bảo tính an toàn (khó đoán hơn math/rand)
		randomIndex, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(charset))),
		)

		// nếu lỗi khi generate random → trả về error ngay
		if err != nil {
			return "", err
		}

		// lấy ký tự tương ứng trong charset và append vào buffer
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}

	// convert buffer thành string và return kết quả cuối cùng
	return strBuilder.String(), nil
}
