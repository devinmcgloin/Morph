package account

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/views/common"
)

func UserLoginView(w http.ResponseWriter, r *http.Request) error {
	_, valid := session.GetUser(r)
	if valid {
		http.Redirect(w, r, "/", 302)
		return nil
	}
	return common.ExecuteTemplate(w, r, "templates/account/loginView.tmpl", nil)
}
