package publicView

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
)

func FeatureImgView(w http.ResponseWriter, r *http.Request) error {

	shortcode := mux.Vars(r)["shortcode"]

	img, err := mongo.GetFeatureSingleImgView(shortcode)
	if err != nil {
		return spriocError.New(err, "Unable to get image", 523)
	}

	usr, valid := session.GetUser(r)
	if valid {
		img.Auth = usr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(img)

	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}
	return nil
}
