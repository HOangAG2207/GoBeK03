package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *healthHandler) CheckHealth(ctx echo.Context) {
	res, err := h.service.CheckHealth()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal Server Error!",
		})
		return

	}

	ctx.JSON(http.StatusOK, res)
}
