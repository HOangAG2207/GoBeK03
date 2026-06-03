package api

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/HOangAG2207/GoBeK03/docs"
	handler "github.com/HOangAG2207/GoBeK03/internal/handler/health"
	urlHandler "github.com/HOangAG2207/GoBeK03/internal/handler/url"
	repository "github.com/HOangAG2207/GoBeK03/internal/repository/health"
	urlRepo "github.com/HOangAG2207/GoBeK03/internal/repository/url"
	service "github.com/HOangAG2207/GoBeK03/internal/service/health"
	urlService "github.com/HOangAG2207/GoBeK03/internal/service/url"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Engine là interface định nghĩa các hành vi của HTTP server
// Giúp dễ mock khi test hoặc thay đổi implementation
type Engine interface {
	Start() error                                     // Hàm dùng để start server
	ServeHTTP(w http.ResponseWriter, r *http.Request) // Cho phép engine hoạt động như http.Handler
}

// struct engine là implementation cụ thể của Engine
type engine struct {
	app         *echo.Echo    // instance của Echo framework (HTTP server)
	config      *Config       // cấu hình ứng dụng
	redisClient *redis.Client // Redis client dùng cho caching / storage
}

// EngineOpts dùng để truyền dependency vào khi khởi tạo Engine (DI - dependency injection)
type EngineOpts struct {
	Cfg         *Config       // config ứng dụng
	RedisClient *redis.Client // Redis client
}

// NewEngine khởi tạo một instance của engine
func NewEngine(opts *EngineOpts) Engine {

	// Tạo instance engine
	e := &engine{
		app:         echo.New(),       // khởi tạo Echo
		config:      opts.Cfg,         // gán config
		redisClient: opts.RedisClient, // gán redis client
	}

	// ===== Middleware =====

	// CORS: cho phép tất cả domain truy cập (dùng "*" -> nên hạn chế trong production)
	e.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Log request HTTP (method, path, status, latency,...)
	e.app.Use(middleware.RequestLogger())

	// Recover: bắt panic để tránh crash server
	e.app.Use(middleware.Recover())

	// ===== Init routes =====
	e.InitRoutes()

	return e
}

// InitRoutes: nơi khai báo toàn bộ routes của ứng dụng
func (e *engine) InitRoutes() {

	// ===== HEALTH CHECK =====

	// Khởi tạo repository (tầng data access)
	checkHealthRepo := repository.NewHealthRepository(e.redisClient)

	// Khởi tạo service (business logic)
	checkHealthService := service.NewHealthService(
		checkHealthRepo,
		e.config.ServiceName,
		e.config.InstanceID,
	)

	// Khởi tạo handler (HTTP layer)
	checkHealthHandler := handler.NewHealthHandler(checkHealthService)

	// ===== SHORTEN URL =====

	// Repository dùng Redis để lưu URL
	shortenUrlRepo := urlRepo.NewUrlRepository(e.redisClient, 0)

	// Service xử lý logic rút gọn URL
	shortenUrlService := urlService.NewUrlService(shortenUrlRepo, 0)

	// Handler nhận request HTTP
	shortenUrlHandler := urlHandler.NewUrlHandler(shortenUrlService)

	// ===== ROUTES =====
	// Redirect từ "/" sang Swagger docs
	// e.app.GET("/", func(c echo.Context) error {
	// 	return c.Redirect(http.StatusFound, "/v1/docs/index.html")
	// })
	// Group tất cả API dưới prefix /api
	apiGroup := e.app.Group("/v1")

	// Swagger UI (docs)
	apiGroup.GET("/docs/*", echoSwagger.WrapHandler)

	// Health check endpoint
	apiGroup.GET("/health-check", checkHealthHandler.CheckHealth)

	// API rút gọn URL
	apiGroup.POST("/links/shorten", shortenUrlHandler.ShortenURL)

}

// Start: chạy HTTP server
func (e *engine) Start() error {

	// ===== 1. Lấy port từ config =====
	port := e.config.AppPort

	// Nếu port chưa có dấu ":" thì thêm vào
	// Echo yêu cầu format dạng ":8080"
	if port[0] != ':' {
		port = fmt.Sprintf(":%s", port)
	}

	// ===== 2. Start server =====
	log.Printf("Server running at %s\n", port)

	// Chạy server (blocking call)
	return e.app.Start(port)
}

// ServeHTTP giúp engine implement interface http.Handler
// Cho phép dùng engine trong test hoặc integrate với server khác
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}
