package authentication

import (
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/spriocError"
)

var mongo = store.ConnectStore()

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", http.StatusNotImplemented)

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", http.StatusNotImplemented)

}
