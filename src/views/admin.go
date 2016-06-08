package views

import (
	"fmt"
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

	page := ps.ByName("page")
	path := fmt.Sprintf("views/morph/%s.tmpl", page)

	t, err := StandardTemplate(path)
	if err != nil {
		log.Printf("Error while parsing template %s", err)
		http.Redirect(w, r, "/morph", 301)
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

	log.Printf("Entering login display")

	if auth.LoggedIn(r) {
		http.Redirect(w, r, "/morph/dashboard", 301)
	}

	path := "views/morph/login.tmpl"

	t, err := StandardTemplate(path)
	if err != nil {
		log.Printf("Error while parsing template %s", err)
		NotFound(w, r)
		return
	}
	log.Print(auth.LoggedIn(r))

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
