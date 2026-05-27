package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CheckHealth handles HTTP GET requests to the /health-check endpoint and returns the service health status.
//
// @Summary      Check health
// @Description  Returns service health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  model.Health
// @Failure      500  {object}  map[string]string
// @Router       /health-check [get]
func (h *healthHandler) CheckHealth(ctx echo.Context) error {
	res, err := h.service.CheckHealth()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error!",
		})
		return err

	}

	return ctx.JSON(http.StatusOK, res)
}
