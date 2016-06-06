package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/dbase"
	"github.com/julienschmidt/httprouter"
)

// IndexHandler handles the index page which is a grid of pictures
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print("Index Handler")
	log.Printf("url_p = %s", r.URL.Path)
	log.Printf("url   = %s", r.URL.Query())
	title := r.URL.Path
	log.Printf("page_t= %s", title)

	page := dbase.GetAllImgs(dbase.DB)

	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, page)
}

// PictureHandler handles the page for individual pictures.
func PictureHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print("Picture Handler")
	log.Printf("url_p = %s", r.URL.Path)
	title := ps.ByName("p_id")
	Img := dbase.GetImg(title, dbase.DB)
	log.Printf("page_t= %s", title)

	t, err := template.ParseFiles("views/content/photo.html")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, Img)
	if err != nil {
		panic(err)
	}
}
