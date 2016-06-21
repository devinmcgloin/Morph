package account

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
)

func ImageEditorView(w http.ResponseWriter, r *http.Request) error {
	usr, _ := session.GetUser(r)

	images, err := mongo.GetUserProfileView(usr.UserName)
	if err != nil {

		return morphError.New(err, "Unable to get user profile view", 523)
	}

	images.Auth = usr

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(images)

	if err != nil {
		return morphError.New(err, "Unable to write JSON", 523)
	}
	return nil

}
