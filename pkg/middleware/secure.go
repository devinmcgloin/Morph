package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/context"
	"github.com/sprioc/sprioc-core/pkg/authentication"
	"github.com/sprioc/sprioc-core/pkg/handlers"
)

func Secure(f func(http.ResponseWriter, *http.Request) handlers.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO this should be formatted in json
		user, err := authentication.CheckUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		context.Set(r, "auth", user)

		w.Header().Set("Content-Type", "application/json")

		resp := f(w, r)
		w.WriteHeader(resp.Code)

		dat, _ := json.MarshalIndent(resp.Data, "", "    ")
		if resp.Data != nil {
			w.Write(dat) // TODO this writes null if the resp.Data is null.
		} else if resp.Message != "" {
			w.Write(resp.Format())
		}
	}
}

func Unsecure(f func(http.ResponseWriter, *http.Request) handlers.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ip, port, _ := net.SplitHostPort(r.RemoteAddr)
		log.Println(ip, port)

		w.Header().Set("Content-Type", "application/json")

		resp := f(w, r)
		w.WriteHeader(resp.Code)

		log.Printf("%+v", resp)

		dat, _ := json.MarshalIndent(resp.Data, "", "    ")
		if resp.Data != nil {
			w.Write(dat) // TODO this writes null if the resp.Data is null.
		} else if resp.Message != "" {
			w.Write(resp.Format())
		}
	}
}
