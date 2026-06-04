package url

import (
	"errors"
	"net/http"

	urlService "github.com/HOangAG2207/GoBeK03/internal/service/url"
	"github.com/labstack/echo/v4"
)

// RedirectURL godoc
// @Summary Redirect short URL
// @Description Redirect to original URL by short code
// @Tags links
// @Accept json
// @Produce json
// @Param code path string true "Short URL code"
// @Success 302 {string} string "Redirect to original URL"
// @Failure 400 {object} map[string]string "Bad request (invalid code or not found)"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/links/redirect/{code} [get]
func (h *urlHandler) RedirectURL(ctx echo.Context) error {
	code := ctx.Param("code")
	if code == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	url, err := h.service.GetURL(ctx.Request().Context(), code)
	if err != nil {
		if errors.Is(err, urlService.ErrCodeNotFound) {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "URL not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return ctx.Redirect(http.StatusFound, url)
}
