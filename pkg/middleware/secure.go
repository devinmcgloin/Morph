package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func Secure(f func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		user, resp := core.CheckUser(r)
		if !resp.Ok() {
			log.Println(resp)
			w.WriteHeader(resp.Code)
			w.Write(resp.Format())
			return
		}

		context.Set(r, "auth", user)

		setIP(r)

		resp = f(w, r)

		w.WriteHeader(resp.Code)

		dat, err := JSONMarshal(resp.Data, true)
		if err != nil {
			log.Println(err)
		}

		if resp.Data != nil {
			w.Write(dat) // TODO this writes null if the resp.Data is null.
		} else if resp.Message != "" {
			w.Write(resp.Format())
		}
	}
}

func Unsecure(f func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		setIP(r)

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

func JSONMarshal(v interface{}, unescape bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "    ")

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func setIP(r *http.Request) {
	ips, ok := r.Header["x-forwarded-for"]
	if !ok {
		log.Println(ips, ok)
	}

	log.Println(strings.Join(ips, "///"))

	trueIP := ips[len(ips)]

	context.Set(r, "ip", trueIP)

}
