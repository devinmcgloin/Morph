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
		log.Fatalf("Redirecting from %s", page)

	}

	t.Execute(w, page)
}
