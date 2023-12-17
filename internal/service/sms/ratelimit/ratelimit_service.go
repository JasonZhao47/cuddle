package ratelimit

import (
	"context"
	"errors"
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"github.com/jasonzhao47/cuddle/pkg/ginx/middleware/ratelimit"
)

type RateLimitSMSServiceV1 struct {
	sms.Service
	limiter ratelimit.Limiter
	key     string
}

func NewRateLimitSMSServiceV1(service sms.Service, limiter ratelimit.Limiter, key string) sms.Service {
	return &RateLimitSMSServiceV1{
		Service: service,
		limiter: limiter,
		key:     key,
	}
}

func (r *RateLimitSMSServiceV1) Send(ctx context.Context, tplId string, args []string, phoneNums []string) error {
	limited, err := r.limiter.Limit(ctx, "")
	if err != nil {
		return err
	}
	if limited {
		return errors.New("触发了限流")
	}
	return r.Service.Send(ctx, tplId, args, phoneNums)
}
