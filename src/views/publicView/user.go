package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
)

func UserProfileView(w http.ResponseWriter, r *http.Request) error {

	UserName := mux.Vars(r)["username"]

	user, err := mongo.GetUserProfileView(UserName)
	if err != nil {
		return morphError.New(err, "Unable to fetch user", 523)
	}

	usr, valid := session.GetUser(r)
	if valid {
		user.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/public/userView.tmpl", user)

}
