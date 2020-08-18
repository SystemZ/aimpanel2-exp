package router

import (
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

var Limiter throttled.HTTPRateLimiter

func InitRateLimit() error {
	store, err := memstore.New(65536)
	if err != nil {
		return err
	}

	quota := throttled.RateQuota{
		MaxBurst: 360,
		MaxRate:  throttled.PerMin(180),
	}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return err
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{RemoteAddr: true},
	}

	Limiter = httpRateLimiter

	return nil
}
