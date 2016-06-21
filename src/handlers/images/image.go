package images

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func ImageHandler(w http.ResponseWriter, r *http.Request) error {
	shortcode := mux.Vars(r)["shortcode"]

	var image model.Image

	err := r.ParseForm()
	if err != nil {

		return spriocError.New(err, "Unable to Parse Form", 523)

	}

	image, err = mongo.GetImageByShortCode(shortcode)
	if err != nil {
		return spriocError.New(err, "Image Not Found", 404)
	}

	err = decoder.Decode(&image, r.PostForm)

	if err != nil {
		return spriocError.New(err, "Unable to decode form", 523)
	}

	mongo.UpdateImage(image)

	// TODO format this so it goes back to the exact right place.
	//	newURL := fmt.Sprintf("/account/images/", shortcode)

	http.Redirect(w, r, "/account/images/", 302)

	return nil
}
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
