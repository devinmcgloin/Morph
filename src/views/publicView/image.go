package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
)

func FeatureImgView(w http.ResponseWriter, r *http.Request) error {

	shortcode := mux.Vars(r)["shortcode"]

	img, err := mongo.GetFeatureSingleImgView(shortcode)
	if err != nil {
		return morphError.New(err, "Unable to get image", 523)
	}

	usr, valid := session.GetUser(r)
	if valid {
		img.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/public/imageView.tmpl", img)
}
