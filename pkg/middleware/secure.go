package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/rsp"
)

// TODO remove access crontrol for production

func Secure(f func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

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
			w.Write(dat) // TODO this writes null if the resp.Data is nil.
		} else if resp.Message != "" {
			w.Write(resp.Format())
		}

	}
}

func Unsecure(f func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

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

	ips, ok := r.Header["X-Forwarded-For"]
	if !ok {
		return
	}

	trueIP := ips[len(ips)-1]

	log.Println(trueIP)

	context.Set(r, "ip", trueIP)
}
