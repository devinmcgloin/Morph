package handler

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/fokal/fokal-core/pkg/domain"

	"bytes"
	"encoding/json"

	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/context"
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
	if rsp.Data == nil {
		return []byte("")
	}
	b, _ := json.MarshalIndent(rsp.Data, "", "    ")

	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	return b
}

// A (simple) example of our application-wide configuration.
type State struct {
	Local bool
	Port  int

	StorageService    StorageState
	AuthService       domain.AuthenticationService
	CacheService      domain.CacheService
	ColorService      domain.ColorService
	PermissionService domain.PermissionService
	SearchService     domain.SearchService
	StreamService     domain.StreamService
	TagService        domain.TagService
	UserService       domain.UserService
	VisionService     domain.VisionService
	ImageService      domain.ImageService
}
type StorageState struct {
	Avatar  domain.StorageService
	Content domain.StorageService
}

// Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	*State
	H func(e *State, w http.ResponseWriter, r *http.Request) (*Response, error)
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.H(h.State, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			if e.Status() >= 500 {
				log.Info("Capturing raven error")
				raven.CaptureError(err, RavenTags(h.State, r))
			}
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			w.WriteHeader(e.Status())
			j, _ := json.Marshal(map[string]interface{}{
				"code": e.Status(),
				"err":  e.Error(),
			})

			w.Write(j)
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			log.Printf("Generating Tags: %+v", RavenTags(h.State, r))
			raven.CaptureError(err, RavenTags(h.State, r))
			log.Printf("HTTP %d - %s", http.StatusInternalServerError, e.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(res.Code)
		w.Write(res.Format())
	}

}

func Options(opts ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Inside Options")
		w.Header().Set("Allow", strings.Join(opts, ", "))
		w.WriteHeader(http.StatusOK)
	})
}

func RavenTags(h *State, r *http.Request) map[string]string {
	tags := map[string]string{}

	if h.Local {
		tags["environment"] = "development"
	} else {
		tags["environment"] = "production"
	}

	contextTags := []string{"ip", "uuid"}
	for _, t := range contextTags {
		value, ok := context.GetOk(r, t)
		if ok {
			tags[t] = fmt.Sprintf("%+v", value)
		}
	}
	return tags
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	j, _ := json.Marshal(map[string]interface{}{
		"code": 404,
		"err":  "Endpoint does not exist.",
	})

	w.Write(j)
}

func StatusHandler(store *State, w http.ResponseWriter, r *http.Request) (Response, error) {
	return Response{Code: http.StatusOK}, nil
}
