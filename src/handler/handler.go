package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/schema"
	"github.com/julienschmidt/httprouter"
)

// IndexHandler handles the index page which is a grid of pictures
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print("Index Handler")
	log.Printf("url_p = %s", r.URL.Path)
	log.Printf("url   = %s", r.URL.Query())
	title := r.URL.Path
	log.Printf("page_t= %s", title)

	page := schema.ImgCollection{
		Title:  "Index",
		NumImg: 10}

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
	title := ps.ByName("img")
	image := schema.ImgPage{
		Title:     title,
		ImgURL:    "/content/beegden-the_netherlands-55.jpg",
		Desc:      "The Netherlands",
		PhotoMeta: schema.PhotoMeta{FStop: 20, ShutterSpeed: 250, FOV: 12, ISO: 200}}
	log.Printf("page_t= %s", title)

	t, err := template.ParseFiles("views/content/photo.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, image)
}
