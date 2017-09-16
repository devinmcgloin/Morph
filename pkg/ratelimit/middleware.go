package ratelimit

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
)

func RateLimit(h http.Handler) http.Handler {
	return tollbooth.LimitHandler(tollbooth.NewLimiter(1, time.Second*30, nil), h)
}
