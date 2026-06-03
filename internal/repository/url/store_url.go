package url

import (
	"context"
	"time"
)

func (r *urlRepository) StoreURL(ctx context.Context, code, url string, exp int64) error {
	if exp <= 0 {
		exp = int64(r.exptime)
	}
	return r.redis.Set(ctx, code, url, time.Duration(exp)*time.Second).Err()
}
