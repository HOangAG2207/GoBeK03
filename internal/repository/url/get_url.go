package url

import "context"

func (r *urlRepository) GetURL(ctx context.Context, code string) (string, error) {
	return r.redis.Get(ctx, code).Result()
}
