package url

import (
	"context"
	"errors"

	"github.com/HOangAG2207/GoBeK03/internal/repository/url"
)

const (
	defaultUrlLengthCode = 7
	maxRetryAttempts     = 10
)

var (
	ErrCodeNotFound       = errors.New("shortened URL not found")
	ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded for generating unique URL code")
)

//go:generate mockery --name UrlService --filename url_service_mock.go --output ./mocks

type UrlService interface {
	ShortenURL(ctx context.Context, url string, exp int64) (string, error)
}
type urlService struct {
	repo          url.UrlRepository
	urlLengthCode int
}

func NewUrlService(repo url.UrlRepository, urlLengthCode int) UrlService {
	if urlLengthCode == 0 {
		urlLengthCode = defaultUrlLengthCode
	}
	return &urlService{
		repo:          repo,
		urlLengthCode: urlLengthCode,
	}
}
