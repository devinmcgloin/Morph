package account

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/auth"
)

func UserLogoutView(w http.ResponseWriter, r *http.Request) error {
	w = auth.LogoutUser(w, r)
	return nil

}
