package ratelimit

import "context"

// NewLimiter 理论上来说应该是有不同实现的，所以分开写

//go:generate mockgen -source=pkg/ginx/middleware/ratelimit/limiter.go -package=limitermocks -destination=./mocks/limiter.mock.go Limiter
type Limiter interface {
	Limit(ctx context.Context, biz string) (bool, error)
}
