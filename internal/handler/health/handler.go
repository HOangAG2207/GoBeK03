package handler

import (
	service "github.com/HOangAG2207/GoBeK03/internal/service/health"
	"github.com/labstack/echo/v4"
)

type Health interface {
	CheckHealth(ctx echo.Context)
}
type healthHandler struct {
	service service.Health
}

func NewHealth(svc service.Health) Health {
	return &healthHandler{
		service: svc,
	}
}
