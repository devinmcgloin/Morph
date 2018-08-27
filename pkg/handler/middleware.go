package handler

import (
	"net/http"
	"time"

	"github.com/newrelic/go-agent"

	"github.com/getsentry/raven-go"
)

type Middleware struct {
	*State
	M func(state *State, next http.Handler) http.Handler
}

func (m Middleware) Handler(next http.Handler) http.Handler {
	return m.M(m.State, next)
}

func Timeout(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, time.Minute, "Application has timed out.")
}

func SentryRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(
		raven.RecoveryHandler(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}))
}

func NewRelic(app newrelic.Application) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		_, middleware := newrelic.WrapHandleFunc(app, "/*", h.ServeHTTP)
		return http.HandlerFunc(middleware)

	}
}
