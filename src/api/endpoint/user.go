package endpoint

import (
	"net/http"

	"github.com/devinmcgloin/sprioc/src/spriocError"
)

func UserHandler(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not implemented", 404)
}
