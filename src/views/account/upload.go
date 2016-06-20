package account

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
)

func UploadView(w http.ResponseWriter, r *http.Request) error {

	var view model.DefaultView

	usr, valid := session.GetUser(r)
	if valid {
		view.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/account/uploadView.tmpl", view)
}
