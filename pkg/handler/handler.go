package handler

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/fokal/fokal-core/pkg/log"

	"github.com/fokal/fokal-core/pkg/services/authentication"
	"github.com/fokal/fokal-core/pkg/services/cache"
	"github.com/fokal/fokal-core/pkg/services/color"
	"github.com/fokal/fokal-core/pkg/services/image"
	"github.com/fokal/fokal-core/pkg/services/permission"
	"github.com/fokal/fokal-core/pkg/services/search"
	"github.com/fokal/fokal-core/pkg/services/storage"
	"github.com/fokal/fokal-core/pkg/services/stream"
	"github.com/fokal/fokal-core/pkg/services/tag"
	"github.com/fokal/fokal-core/pkg/services/user"
	"github.com/fokal/fokal-core/pkg/services/vision"

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
	AuthService       authentication.Service
	CacheService      cache.Service
	ColorService      color.Service
	PermissionService permission.Service
	SearchService     search.Service
	StreamService     stream.Service
	TagService        tag.Service
	UserService       user.Service
	VisionService     vision.Service
	ImageService      image.Service
}
type StorageState struct {
	Avatar  storage.Service
	Content storage.Service
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
	ctx := r.Context()
	if err != nil {
		switch e := err.(type) {
		case Error:
			if e.Status() >= 500 {
				logrus.Info("Capturing raven error")
				raven.CaptureError(err, RavenTags(h.State, r))
			}
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.WithContext(ctx).WithFields(logrus.Fields{
				"status-code": e.Status(),
				"error":       e,
			}).Infof("unable to fulfill request")
			w.WriteHeader(e.Status())
			j, _ := json.Marshal(map[string]interface{}{
				"code": e.Status(),
				"err":  e.Error(),
			})

			w.Write(j)
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			logrus.Printf("Generating Tags: %+v", RavenTags(h.State, r))
			raven.CaptureError(err, RavenTags(h.State, r))
			log.WithContext(ctx).WithFields(logrus.Fields{
				"status-code": http.StatusInternalServerError,
				"error":       e,
			}).Infof("unable to fulfill request")
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
		logrus.Println("Inside Options")
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
