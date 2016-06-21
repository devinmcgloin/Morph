package account

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
)

func UploadView(w http.ResponseWriter, r *http.Request) error {

	var view model.DefaultView

	usr, valid := session.GetUser(r)
	if valid {
		view.Auth = usr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(usr)

	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}
	return nil
}
