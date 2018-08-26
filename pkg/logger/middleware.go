package logger

import (
	"context"
	"net/http"

	"net"
	"strings"

	"github.com/satori/go.uuid"
)

const (
	requestIDKey = iota
	ipIDKey
)

func UUID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIDKey, uuid.NewV4())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
			addresses := strings.Split(r.Header.Get(h), ",")
			// march from right to left until we get a public address
			// that will be the address right before our proxy.
			for i := len(addresses) - 1; i >= 0; i-- {
				ip := strings.TrimSpace(addresses[i])
				// header can contain spaces too, strip those out.
				realIP := net.ParseIP(ip)
				if !realIP.IsGlobalUnicast() {
					// bad address, go to next
					continue
				}
				ctx = context.WithValue(ctx, ipIDKey, realIP.String())
				break
			}
		}

		h.ServeHTTP(w, r.WithContext(ctx))

	})
}

func ContentTypeJSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
