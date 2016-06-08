package handler

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/dbase"
	"github.com/devinmcgloin/morph/src/schema"
	"github.com/julienschmidt/httprouter"
)

// IndexHandler handles the index page which is a grid of pictures
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	collection, err := dbase.GetAllImgs()
	if err != nil {
		log.Printf("Error while getting all images from DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	renderCollection(w, collection)

}

// PictureHandler handles the page for individual pictures.
func PictureHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	title := ps.ByName("i_id")
	Img, err := dbase.GetImg(title)
	if err != nil {
		log.Printf("Error while getting images from DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	t, err := StandardTemplate("views/content/image.tmpl")
	if err != nil {
		log.Printf("Error while parsing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = t.Execute(w, Img)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func CategoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	album := ps.ByName("album")
	collection, err := dbase.GetAlbum(album)
	if err != nil {
		log.Printf("Error while getting all images from DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	renderCollection(w, collection)
}

func renderCollection(w http.ResponseWriter, collection schema.ImgCollection) {

	t, err := StandardTemplate("views/content/collection.tmpl")
	if err != nil {
		log.Printf("Error while parsing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = t.Execute(w, collection)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
