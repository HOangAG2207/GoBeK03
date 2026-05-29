package api

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/HOangAG2207/GoBeK03/docs"
	handler "github.com/HOangAG2207/GoBeK03/internal/handler/health"
	repository "github.com/HOangAG2207/GoBeK03/internal/repository/health"
	service "github.com/HOangAG2207/GoBeK03/internal/service/health"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
type engine struct {
	app         *echo.Echo
	config      *Config
	redisClient *redis.Client
}

type EngineOpts struct {
	Cfg         *Config
	RedisClient *redis.Client
}

func NewEngine(opts *EngineOpts) Engine {
	e := &engine{
		app:         echo.New(),
		config:      opts.Cfg,
		redisClient: opts.RedisClient,
	}
	e.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// Log incoming HTTP requests
	e.app.Use(middleware.RequestLogger())

	// Recover from panics and prevent server crash
	e.app.Use(middleware.Recover())

	//Init Routes
	e.InitRoutes()

	return e
}

// Init Routes
func (e *engine) InitRoutes() {

	//health-check
	checkHealthRepo := repository.NewHealth()
	checkHealthService := service.NewHealth(checkHealthRepo, e.config.ServiceName, e.config.InstanceID)
	checkHealthHandler := handler.NewHealth(checkHealthService)

	e.app.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
	})

	apiGroup := e.app.Group("/api")
	apiGroup.GET("/docs/*", echoSwagger.WrapHandler)

	apiGroup.GET("/health-check", checkHealthHandler.CheckHealth)
}

// Start runs the HTTP server
func (e *engine) Start() error {
	// ===== 1. Define server port =====
	port := e.config.AppPort
	// port := ":8080"
	// Ensure port has ":" prefix (required by Echo)
	if port[0] != ':' {
		port = fmt.Sprintf(":%s", port)
	}
	// ===== 2. Start server =====
	log.Printf("Server running at %s\n", port)
	// Start Echo server
	return e.app.Start(port)
}
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}
