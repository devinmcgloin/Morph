package endpoint

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	shortcode := mux.Vars(r)["shortcode"]

	var image model.Image

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	image, err = mongo.GetImageByShortCode(shortcode)
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}

	err = decoder.Decode(&image, r.PostForm)

	if err != nil {
		log.Println(err)
		common.SomethingsWrong(w, r, err)
		return
	}

	mongo.UpdateImage(image)

	newURL := fmt.Sprintf("/i/%s/edit", shortcode)

	http.Redirect(w, r, newURL, 302)

}
