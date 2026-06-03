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
// @Router       /v1/health-check [get]
func (h *healthHandler) CheckHealth(ctx echo.Context) error {

	// ===== 1. Gọi service layer =====

	// Gọi business logic để lấy thông tin health
	res, err := h.service.CheckHealth(ctx.Request().Context())

	// ===== 2. Handle error =====
	if err != nil {

		// Trả về HTTP 500 + JSON message lỗi
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error!",
		})

		// Trả error cho Echo (có thể được middleware log lại)
		return err
	}

	// ===== 3. Return success =====

	// Trả về HTTP 200 + dữ liệu health
	return ctx.JSON(http.StatusOK, res)
}
