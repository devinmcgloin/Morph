package common

import (
	"log"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	RenderStatic(w, r, "templates/pages/404.tmpl")
}

func SomethingsWrong(w http.ResponseWriter, r *http.Request, outsideError error) {
	log.Println("Something's wrong!")

	t, err := StandardTemplate("templates/pages/internalError.tmpl")
	if err != nil {
		log.Printf("Error while parsing template: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = t.Execute(w, outsideError)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
