package url

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (s *urlService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	url, err := s.repo.GetURL(ctx, code)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}
		return "", fmt.Errorf("get original url: %w", err)
	}
	return url, nil
}
