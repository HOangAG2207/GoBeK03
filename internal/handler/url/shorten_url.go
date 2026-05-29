package url

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type shortenURLRequest struct {
	URL string `json:"url" example:"https://www.google.com"`
	// ExpireIn int    `json:"exp" binding:"required,lte=604800" example:"10000"`
}

type shortenURLResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// ShortenURL godoc
// @Summary      Shorten a URL
// @Description  Generate a shortened URL code from a given URL
// @Tags         url
// @Accept       json
// @Produce      json
// @Param        request  body      shortenURLRequest  true  "URL to shorten"
// @Success      200      {object}  shortenURLResponse
// @Failure      400      {object}  shortenURLResponse
// @Failure      500      {object}  shortenURLResponse
// @Router       /url/shorten [post]
func (h *urlHandler) ShortenURL(ctx echo.Context) error {
	input := new(shortenURLRequest)

	if err := ctx.Bind(input); err != nil || input.URL == "" {
		return ctx.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: InValidRequestPayload.Error(),
		})
	}
	code, err := h.service.ShortenURL(ctx.Request().Context(), input.URL)
	if err != nil {
		log.Error().
			Str("url", input.URL).
			Err(err).
			Msg("service return error when shorten url")

		return ctx.JSON(http.StatusInternalServerError, shortenURLResponse{
			// Message: InternalServerError.Error(),
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, shortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
