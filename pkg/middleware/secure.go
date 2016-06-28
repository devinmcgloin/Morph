package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/context"
	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func Secure(f func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO this should be formatted in json
		user, err := core.CheckUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		context.Set(r, "auth", user)

		w.Header().Set("Content-Type", "application/json")

		resp := f(w, r)
		w.WriteHeader(resp.Code)

		dat, _ := JSONMarshal(resp.Data, true)
		if resp.Data != nil {
			w.Write(dat) // TODO this writes null if the resp.Data is null.
		} else if resp.Message != "" {
			w.Write(resp.Format())
		}
	}
}

func Unsecure(f func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ip, port, _ := net.SplitHostPort(r.RemoteAddr)
		log.Println(ip, port)

		w.Header().Set("Content-Type", "application/json")

		resp := f(w, r)
		w.WriteHeader(resp.Code)

		log.Printf("%+v", resp)

		dat, _ := JSONMarshal(resp.Data, true)
		if resp.Data != nil {
			w.Write(dat) // TODO this writes null if the resp.Data is null.
		} else if resp.Message != "" {
			w.Write(resp.Format())
		}
	}
}

func JSONMarshal(v interface{}, unescape bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "    ")

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
