package url

import "context"

func (r *urlRepository) StoreURLIfAbsent(ctx context.Context, code, url string) (bool, error) {
	return r.redis.SetNX(ctx, code, url, r.exptime).Result()
}
