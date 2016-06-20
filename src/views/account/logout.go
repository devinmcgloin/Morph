package account

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/auth"
	"github.com/devinmcgloin/morph/src/views/common"
)

func UserLogoutView(w http.ResponseWriter, r *http.Request) error {
	w = auth.LogoutUser(w, r)
	return common.ExecuteTemplate(w, r, "templates/account/logoutView.tmpl", nil)
}
