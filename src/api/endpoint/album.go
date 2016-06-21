package endpoint

import (
	"net/http"

	"github.com/devinmcgloin/sprioc/src/spriocError"
)

func AlbumHandler(w http.ResponseWriter, r *http.Request) error {

	return spriocError.New(nil, "Not Implemented", 404)
}
