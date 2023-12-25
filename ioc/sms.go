package ioc

import (
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"github.com/jasonzhao47/cuddle/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter())
	return localsms.NewService()
}
