package ratelimit

import "context"

type Limiter interface {
	// Limit checks whether the request is limited.
	Limit(ctx context.Context, key string) (bool, error)
}
