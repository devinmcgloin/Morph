package collections

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
)

func MostRecentView(w http.ResponseWriter, r *http.Request) error {

	var images model.CollectionView

	images, err := mongo.GetNumMostRecentImg(10)
	if err != nil {
		return spriocError.New(err, "Unable to get collection", 523)

	}

	usr, valid := session.GetUser(r)
	if valid {
		images.Auth = usr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(images)

	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}
	return nil

}
