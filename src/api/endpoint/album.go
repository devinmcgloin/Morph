package endpoint

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/morphError"
)

func AlbumHandler(w http.ResponseWriter, r *http.Request) error {

	return morphError.New(nil, "Not Implemented", 404)
}
