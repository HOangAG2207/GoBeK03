package url

import (
	"context"
	"errors"

	"github.com/HOangAG2207/GoBeK03/internal/repository/url"
)

const (
	maxRetryAttempts = 10
)

var (
	ErrCodeNotFound       = errors.New("shortened URL not found")
	ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded for generating unique URL code")
)

type UrlService interface {
	ShortenURL(ctx context.Context, url string) (string, error)
}
type urlService struct {
	repo          url.UrlRepository
	urlLengthCode int
}

func NewUrlService(repo url.UrlRepository, urlLengthCode int) UrlService {
	if urlLengthCode == 0 {
		urlLengthCode = 10
	}
	return &urlService{
		repo:          repo,
		urlLengthCode: urlLengthCode,
	}
}
