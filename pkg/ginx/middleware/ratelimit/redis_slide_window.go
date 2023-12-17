package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var slideWindow string

type WindowLimiter struct {
	client   redis.Cmdable
	interval time.Duration
	rate     int
}

func (s *WindowLimiter) Limit(ctx context.Context, biz string) (bool, error) {
	return s.client.Eval(ctx, slideWindow, []string{s.key(biz)}, s.interval.Milliseconds(), s.rate, time.Now().UnixMilli()).Bool()
}

func NewSlidingWindowLimiter(client redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &WindowLimiter{
		client:   client,
		interval: interval,
		rate:     rate,
	}
}

func (s *WindowLimiter) key(biz string) string {
	return fmt.Sprintf("limiter:%s", biz)
}
