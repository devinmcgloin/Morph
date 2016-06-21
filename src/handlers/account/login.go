package account

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
)

func UserLoginView(w http.ResponseWriter, r *http.Request) error {
	_, valid := session.GetUser(r)
	if valid {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(valid)

	if err != nil {
		return morphError.New(err, "Unable to write JSON", 523)
	}
	return nil
}
