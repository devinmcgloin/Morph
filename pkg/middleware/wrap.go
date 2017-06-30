package middleware

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/rsp"
)

func Wrap(f ...func(http.ResponseWriter, *http.Request) rsp.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var rsp rsp.Response
		for _, fun := range f {
			rsp = fun(w, r)
			if !rsp.Ok() {
				w.WriteHeader(rsp.Code)
				w.Write(rsp.Format())
			}
		}
		w.WriteHeader(rsp.Code)

		dat, err := JSONMarshal(rsp.Data, true)
		if err != nil {
			log.Println(err)
		}

		if rsp.Data != nil {
			w.Write(dat) // TODO this writes null if the resp.Data is nil.
		} else if rsp.Message != "" {
			w.Write(rsp.Format())
		}
	}
}
