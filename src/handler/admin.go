package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/auth"
	"github.com/julienschmidt/httprouter"
)

func AdminHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	log.Print("entering admin handler")

	if !auth.LoggedIn(r) {
		log.Printf("Not Logged in, Redirecting.")
		http.Redirect(w, r, "/morph", 301)
		return
	}

	page := ps.ByName("type")
	path := fmt.Sprintf("views/morph/%s.html", page)

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Error while parsing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = t.Execute(w, page)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func LoginDisplay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := "views/morph/login.html"

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("Error while parsing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
