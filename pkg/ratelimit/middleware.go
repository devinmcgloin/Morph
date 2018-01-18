package ratelimit

import (
	"net/http"

	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

func RateLimit(h http.Handler) http.Handler {
	return tollbooth.LimitHandler(tollbooth.NewLimiter(5, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	}), h)
}
