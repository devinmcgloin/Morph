package publicView

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
)

func UserProfileView(w http.ResponseWriter, r *http.Request) error {

	UserName := mux.Vars(r)["username"]

	user, err := mongo.GetUserProfileView(UserName)
	if err != nil {
		return spriocError.New(err, "Unable to fetch user", 523)
	}

	usr, valid := session.GetUser(r)
	if valid {
		user.Auth = usr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(usr)

	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}
	return nil
}
