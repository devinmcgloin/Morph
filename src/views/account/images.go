package account

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
	"github.com/devinmcgloin/morph/src/views/common"
)

func ImageEditorView(w http.ResponseWriter, r *http.Request) error {
	usr, _ := session.GetUser(r)

	images, err := mongo.GetUserProfileView(usr.UserName)
	if err != nil {

		return morphError.New(err, "Unable to get user profile view", 523)
	}

	images.Auth = usr

	return common.ExecuteTemplate(w, r, "templates/account/imagesEditView.tmpl", images)
}
