package endpoint

import (
	"net/http"

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
