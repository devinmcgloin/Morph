package handler

import (
	"html/template"
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
	log.Print("Picture Handler")
	log.Printf("url_p = %s", r.URL.Path)
	title := ps.ByName("i_id")
	Img, err := dbase.GetImg(title)
	if err != nil {
		log.Printf("Error while getting images from DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	log.Printf("page_t= %s", title)

	t, err := template.ParseFiles("views/content/image.html")
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
	category := ps.ByName("category")
	collection, err := dbase.GetCategory(category)
	if err != nil {
		log.Printf("Error while getting all images from DB %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	renderCollection(w, collection)
}

func renderCollection(w http.ResponseWriter, collection schema.ImgCollection) {

	t, err := template.ParseFiles("views/content/collection.html")
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
