package handler

import (
	"log"
	"net/http"

	"bytes"
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
	vision "google.golang.org/api/vision/v1"
	"googlemaps.github.io/maps"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

type Response struct {
	Code int
	Data interface{}
}

func (rsp Response) Format() []byte {
	b, _ := json.MarshalIndent(rsp.Data, "", "    ")

	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	return b
}

// A (simple) example of our application-wide configuration.
type State struct {
	DB *sqlx.DB
	//ES
	RD     *redis.Pool
	Local  bool
	Vision *vision.Service
	Maps   *maps.Client
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	*State
	H func(e *State, w http.ResponseWriter, r *http.Request) (Response, error)
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.H(h.State, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(res.Code)
		w.Write(res.Format())
	}

}
