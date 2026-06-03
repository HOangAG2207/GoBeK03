package api

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config là struct chứa toàn bộ cấu hình của ứng dụng
// Giá trị sẽ được load từ environment variables (.env hoặc system env)
type Config struct {

	// AppPort: port mà server sẽ chạy
	// envconfig:"PORT" → đọc từ biến môi trường APP_PORT (do prefix "APP")
	// default:"8081" → nếu không có env thì dùng 8081
	AppPort string `envconfig:"PORT" default:"8081"`

	// ServiceName: tên service (dùng cho logging, monitoring, tracing...)
	ServiceName string `envconfig:"SERVICE_NAME" default:"GoBe-K03"`

	// InstanceID: định danh instance (hữu ích khi scale nhiều instance)
	// ví dụ: pod-id trong Kubernetes
	InstanceID string `envconfig:"INSTANCE_ID" default:""`
}

// NewConfig: load config từ .env + environment variables
func NewConfig() (*Config, error) {

	// ===== 1. Load file .env =====

	// godotenv.Load() sẽ đọc file .env và set vào environment variables
	// ví dụ: APP_PORT=8080
	if err := godotenv.Load(); err != nil {

		// Nếu không có file .env → KHÔNG lỗi, chỉ log warning
		// Điều này hữu ích khi chạy production (env đã có sẵn)
		log.Println(".env not found")
	}

	// ===== 2. Khởi tạo struct config =====
	var cfg Config

	// ===== 3. Parse env → struct =====

	// envconfig.Process("APP", &cfg)
	// nghĩa là:
	// sẽ đọc các biến môi trường có prefix APP_
	//
	// Ví dụ mapping:
	// APP_PORT=8080          → cfg.AppPort
	// APP_SERVICE_NAME=abc   → cfg.ServiceName
	// APP_INSTANCE_ID=xyz    → cfg.InstanceID
	if err := envconfig.Process("APP", &cfg); err != nil {

		// Nếu parse lỗi (type mismatch, thiếu field bắt buộc...)
		return nil, err
	}

	// ===== 4. Trả về config =====
	return &cfg, nil
}
