package endpoint

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
)

var decoder = schema.NewDecoder()

func ImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	imageShortTitle := ps.ByName("ImageShortTitle")

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	img, err := mongo.GetImageByTitle(imageShortTitle)
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}

	err = decoder.Decode(&img, r.PostForm)

	log.Println(r.PostForm)

	log.Println(img)

	if err != nil {
		log.Println(err)
		common.SomethingsWrong(w, r, err)
		return
	}

	mongo.UpdateImage(img)

	newUrl := fmt.Sprintf("/i/%s/edit", imageShortTitle)

	http.Redirect(w, r, newUrl, 302)

}
