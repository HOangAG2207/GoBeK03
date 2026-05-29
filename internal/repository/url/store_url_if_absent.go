package url

import (
	"context"
	"fmt"
)

func (r *urlRepository) StoreURLIfAbsent(ctx context.Context, code, url string) (bool, error) {
	ok, err := r.redis.SetNX(ctx, code, url, r.exptime).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx failed: %w", err)
	}
	return ok, nil
}
