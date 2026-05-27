package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/HOangAG2207/GoBeK03/internal/config"
	handler "github.com/HOangAG2207/GoBeK03/internal/handler/health"
	repository "github.com/HOangAG2207/GoBeK03/internal/repository/health"
	service "github.com/HOangAG2207/GoBeK03/internal/service/health"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
type engine struct {
	app    *echo.Echo
	config *config.Config
}

type EngineOpts struct {
	Cfg *config.Config
}

func NewEngine(opts *EngineOpts) Engine {
	e := &engine{
		app:    echo.New(),
		config: opts.Cfg,
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

	apiGroup := e.app.Group("/api")
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
