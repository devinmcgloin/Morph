package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// IndexHandler handles the index page which is a grid of pictures
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print("Index Handler")
	log.Printf("url_p = %s", r.URL.Path)
	log.Printf("url   = %s", r.URL.Query())
	title := r.URL.Path
	log.Printf("page_t= %s", title)

	page := ImgCollection{
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
	image := Img{
		Title:     title,
		ImgURL:    "/content/beegden-the_netherlands-55.jpg",
		Desc:      "The Netherlands",
		PhotoMeta: PhotoMeta{FStop: 20, ShutterSpeed: 250, FOV: 12, ISO: 200}}
	log.Printf("page_t= %s", title)

	t, err := template.ParseFiles("views/content/photo.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, image)
}

// Img is the basic type for single images.
type Img struct {
	Title     string
	Desc      string
	ImgURL    string
	PhotoMeta PhotoMeta
}

// PhotoMeta contains all the meta information about a specific image
type PhotoMeta struct {
	FStop        int
	ShutterSpeed int
	FOV          int
	ISO          int
}

// ImgPage contains all the proper information for rendering a single photo
type ImgPage struct {
	Title    string
	Img      Img
	PageMeta PageMeta
}

// ImgCollection includes a title and collection of Images.
type ImgCollection struct {
	Title    string
	NumImg   int
	Images   []Img
	PageMeta PageMeta
}

// PageMeta is a type for the meta tags found at the top of the page.
type PageMeta struct {
	Title         string
	Image         string
	URL           string
	Description   string
	Author        string
	AuthorTwitter string
	Keywords      []string
}
