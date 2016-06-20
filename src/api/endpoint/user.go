package endpoint

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/morphError"
)

func UserHandler(w http.ResponseWriter, r *http.Request) error {
	return morphError.New(nil, "Not implemented", 404)
}
