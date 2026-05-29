package url

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UrlRepository interface {
	StoreURL(ctx context.Context, code, url string) error
	GetURL(ctx context.Context, code string) (string, error)
}
type urlRepository struct {
	redis   *redis.Client
	exptime time.Duration
}

func NewUrlRepository(rc *redis.Client, exptime time.Duration) UrlRepository {
	if exptime <= 0 {
		exptime = 24 * time.Hour
	}
	return &urlRepository{
		redis:   rc,
		exptime: exptime,
	}
}
