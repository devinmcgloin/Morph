package security

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/core"
	"github.com/gorilla/context"
)

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, resp := core.CheckUser(r)
		if !resp.Ok() {
			log.Println(resp)
			w.WriteHeader(resp.Code)
			w.Write(resp.Format())
			return
		}

		context.Set(r, "auth", user)

		h.ServeHTTP(w, r)
	})
}
