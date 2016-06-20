package common

import (
	"log"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	t, err := StandardTemplate("templates/public/404.tmpl")
	if err != nil {
		log.Printf("Error while parsing template: %s", err)
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

func SomethingsWrong(w http.ResponseWriter, r *http.Request, outsideError error) {
	t, err := StandardTemplate("templates/public/internalError.tmpl")
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
