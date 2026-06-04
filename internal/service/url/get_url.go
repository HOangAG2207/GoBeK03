package url

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func (s *urlService) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.repo.GetURL(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}
	return url, nil
}
