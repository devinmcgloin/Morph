package status

import (
	"net/http"

	"github.com/fokal/fokal-core/pkg/handler"
)

func StatusHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	return handler.Response{Code: http.StatusOK}, nil
}
