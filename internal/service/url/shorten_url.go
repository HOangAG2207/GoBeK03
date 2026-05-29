package url

import (
	"context"

	"github.com/HOangAG2207/GoBeK03/internal/utils"
)

func (s *urlService) ShortenURL(ctx context.Context, url string) (string, error) {
	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		urlCode, err := utils.GenerateCode(s.urlLengthCode)

		if err != nil {
			return "", err
		}

		stored, err := s.repo.StoreURLIfAbsent(ctx, urlCode, url)
		if err != nil {
			return "", err
		}
		if stored {
			return urlCode, nil // atomically claimed
		}
	}
	return "", ErrMaxRetriesExceeded
}
