package url

import "context"

func (r *urlRepository) StoreURL(ctx context.Context, code, url string) error {
	return r.redis.Set(ctx, code, url, r.exptime).Err()
}
