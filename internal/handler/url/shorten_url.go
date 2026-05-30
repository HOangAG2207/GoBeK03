package url

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type shortenURLRequest struct {
	URL string `json:"url" validate:"required,url" example:"https://www.google.com"`
	Exp int64  `json:"exp" example:"604800"` // thời gian hết hạn (giây)
}

type shortenURLResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ShortenURL godoc
// @Summary      Shorten a URL
// @Description  Generate a shortened URL code from a given URL with optional expiration time
// @Tags         url
// @Accept       json
// @Produce      json
// @Param        request  body      shortenURLRequest  true  "URL to shorten with expiration time (exp in seconds)"
// @Success      200      {object}  shortenURLResponse
// @Failure      400      {object}  shortenURLResponse
// @Failure      500      {object}  shortenURLResponse
// @Router       /v1/links/shorten [post]
func (h *urlHandler) ShortenURL(ctx echo.Context) error {
	input := new(shortenURLRequest)

	if err := ctx.Bind(input); err != nil || input.URL == "" {
		return ctx.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: InValidRequestPayload.Error(),
		})
	}

	// ✅ validate exp
	if input.Exp < 0 {
		return ctx.JSON(http.StatusBadRequest, shortenURLResponse{
			Message: "expiration must be greater than 0",
		})
	}

	// default exp nếu không truyền
	exp := input.Exp
	if exp == 0 {
		exp = 3600
	}

	code, err := h.service.ShortenURL(
		ctx.Request().Context(),
		input.URL,
		exp,
	)

	if err != nil {
		log.Error().
			Str("url", input.URL).
			Int64("exp", exp).
			Err(err).
			Msg("service return error when shorten url")

		return ctx.JSON(http.StatusInternalServerError, shortenURLResponse{
			Message: InternalServerError.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, shortenURLResponse{
		Code:    code,
		Message: "Shorten URL generated successfully!",
	})
}
