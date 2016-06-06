package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AdminHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Print("Admin Handler")
	page := ps.ByName("type")
	path := fmt.Sprintf("views/admin/%s.html", page)

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
