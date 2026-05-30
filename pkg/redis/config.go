package pkgredis

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// config chứa toàn bộ cấu hình kết nối Redis
// Các giá trị sẽ được lấy từ environment variables hoặc .env file
type config struct {

	// Address: địa chỉ Redis server (host:port)
	Address string `envconfig:"REDIS_ADDR" default:"localhost:6379"`

	// Password: mật khẩu Redis (nếu có enable auth)
	Password string `envconfig:"REDIS_PASSWORD" default:""`

	// DB: Redis database index (0-15 mặc định)
	DB int `envconfig:"REDIS_DB" default:"0"`
}

// newConfig đọc cấu hình Redis từ env + .env file
// envprefix dùng để hỗ trợ prefix khi đọc biến môi trường
func newConfig(envprefix string) (*config, error) {

	// ===== 1. Load file .env =====

	// godotenv.Load() sẽ đọc file .env và inject vào environment variables
	// ví dụ: REDIS_ADDR=localhost:6379
	if err := godotenv.Load(); err != nil {

		// nếu không có file .env thì không lỗi, chỉ log warning
		// vì production thường dùng env thật thay vì file .env
		log.Println(".env not found")
	}

	// ===== 2. Khởi tạo config struct =====
	var cfg config

	// ===== 3. Parse environment variables =====

	// envconfig.Process sẽ map environment variables vào struct
	// theo tag envconfig
	//
	// Ví dụ:
	// REDIS_ADDR=localhost:6379
	// REDIS_PASSWORD=123
	// REDIS_DB=0
	if err := envconfig.Process(envprefix, &cfg); err != nil {

		// nếu parse lỗi → trả error lên caller
		return nil, err
	}

	// ===== 4. Return config =====
	return &cfg, nil
}
