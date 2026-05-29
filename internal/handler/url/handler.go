package url

import (
	"errors"

	"github.com/HOangAG2207/GoBeK03/internal/service/url"
	"github.com/labstack/echo/v4"
)

var (
	ErrCodeNotFound       = errors.New("code not found")
	InValidRequestPayload = errors.New("invalid request payload")
	InternalServerError   = errors.New("internal server error")
)

type UrlHandler interface {
	ShortenURL(ctx echo.Context) error
}
type urlHandler struct {
	service url.UrlService
}

func NewUrlHandler(svc url.UrlService) UrlHandler {
	return &urlHandler{
		service: svc,
	}
}
