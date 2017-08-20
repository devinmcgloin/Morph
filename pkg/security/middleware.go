package security

import (
	"net/http"

	"log"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/gorilla/context"
)

func Authenticate(state *handler.State, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := verifyJWT(state, r)
		if err != nil {
			switch e := err.(type) {
			case handler.Error:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				log.Printf("HTTP %d - %s", e.Status(), e.Error())
				http.Error(w, e.Error(), e.Status())
			default:
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		} else {
			context.Set(r, "auth", user)
			next.ServeHTTP(w, r)
		}

	})
}

func SetAuthenticatedUser(state *handler.State, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := verifyJWT(state, r)
		if err == nil {
			context.Set(r, "auth", user)
		}
		next.ServeHTTP(w, r)
	})
}
