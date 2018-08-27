package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"net"
	"strings"

	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/request"
	raven "github.com/getsentry/raven-go"
	"github.com/satori/go.uuid"
)

type Middleware struct {
	*handler.State
	M func(state *handler.State, next http.Handler) http.Handler
}

func (m Middleware) Handler(next http.Handler) http.Handler {
	return m.M(m.State, next)
}

func UUID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, request.IDKey, uuid.NewV4())
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
				ctx = context.WithValue(ctx, request.IPKey, realIP.String())
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

func Cache(state *handler.State, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if state.Local {
			next.ServeHTTP(w, r)
		} else {
			url := r.URL.String()
			b, err := state.CacheService.Get(url)
			if err != nil {
				c := httptest.NewRecorder()
				next.ServeHTTP(c, r)

				for k, v := range c.Header() {
					w.Header()[k] = v
				}

				content := c.Body.Bytes()
				w.WriteHeader(c.Code)
				w.Write(content)

				if c.Code == http.StatusOK {
					log.Printf("Cache: Setting Handler URL: %s\n", url)
					state.CacheService.Set(url, content)
				}
				return
			}
			log.Printf("Cache: Retrieving Handler URL: %s\n", url)
			w.Write(b)
			w.WriteHeader(http.StatusOK)
			return

		}
	})
}

func SetJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.ExtractToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
		} else {
			tokenStr := strings.Replace(tokenStrings, "Bearer ", "", 1)
			ctx = context.WithValue(ctx, request.JWTID, tokenStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func Authenticate(state *handler.State, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		jwt, ok := ctx.Value(request.JWTID).(string)
		if !ok {
			return
		}
		valid, userID, err := state.AuthService.VerifyToken(ctx, jwt)
		if err != nil {
			switch e := err.(type) {
			case handler.Error:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				log.Printf("HTTP %d - %s", e.Status(), e.Error())
				w.WriteHeader(e.Status())
				j, _ := json.Marshal(map[string]interface{}{
					"code": e.Status(),
					"err":  e.Error(),
				})
				w.Write(j)
			default:
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		} else if !valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			ctx = context.WithValue(ctx, request.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

	})
}

func SetAuthenticatedUser(state *handler.State, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		jwt, ok := ctx.Value(request.JWTID).(string)
		if !ok {
			next.ServeHTTP(w, r)
		}
		valid, userID, err := state.AuthService.VerifyToken(ctx, jwt)
		if !valid || err == nil {
			next.ServeHTTP(w, r)
		} else {
			ctx = context.WithValue(ctx, request.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
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

func RateLimit(h http.Handler) http.Handler {
	return tollbooth.LimitHandler(tollbooth.NewLimiter(5, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Second,
	}), h)
}
