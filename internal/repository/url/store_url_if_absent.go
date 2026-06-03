package url

import (
	"context"
	"fmt"
	"time"
)

func (r *urlRepository) StoreURLIfAbsent(ctx context.Context, code, url string, exp int64) (bool, error) {
	if exp <= 0 {
		exp = int64(r.exptime)
	}
	ok, err := r.redis.SetNX(ctx, code, url, time.Duration(exp)*time.Second).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx failed: %w", err)
	}
	return ok, nil
}
