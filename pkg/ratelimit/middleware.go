package ratelimit

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
)

func RateLimit(h http.Handler) http.Handler {
	return tollbooth.LimitHandler(tollbooth.NewLimiter(5, time.Second, nil), h)
}
