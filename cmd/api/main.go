package main

import (
	"log"

	"github.com/HOangAG2207/GoBeK03/internal/api"
	pkgredis "github.com/HOangAG2207/GoBeK03/pkg/redis"
)

// Khai báo metadata cho Swagger (dùng để generate tài liệu API tự động)

// @title           GoBe K03 project API
// Tiêu đề của API, sẽ hiển thị trên giao diện Swagger UI

// @version         1.0
// Phiên bản API hiện tại

// @description     API for GoBe K03
// Mô tả ngắn gọn về chức năng của API

// @host            localhost:8080
// Host mà Swagger sẽ gọi API (thường dùng cho môi trường local/dev)

// @BasePath 		/
// Base path (prefix) cho tất cả các endpoint API

func main() {

	// Gọi hàm NewConfig để load cấu hình (config)
	// Thường config sẽ lấy từ file .env, environment variables hoặc file yaml/json
	cfg, err := api.NewConfig()

	// Kiểm tra nếu có lỗi khi load config
	if err != nil {
		// log.Fatalf sẽ in log và terminate chương trình ngay lập tức
		log.Fatalf("failed to load config: %v", err)
	}

	// Khởi tạo Redis client
	// Tham số "" có thể là địa chỉ Redis (ở đây đang để rỗng -> dùng default trong hàm)
	redisClient, err := pkgredis.NewRedisClient("")

	// Nếu có lỗi khi kết nối Redis
	if err != nil {
		// panic sẽ dừng chương trình và in stack trace (thường dùng cho lỗi nghiêm trọng)
		panic(err)
	}

	// Khởi tạo HTTP engine (Gin hoặc framework bạn đang dùng)
	// Truyền vào các options cần thiết cho app
	app := api.NewEngine(&api.EngineOpts{
		Cfg:         cfg,         // truyền config đã load vào engine
		RedisClient: redisClient, // truyền Redis client để các layer khác sử dụng
	})

	// Start server (chạy HTTP server)
	if err := app.Start(); err != nil {
		// Nếu server không start được thì panic
		panic(err)
	}
}
