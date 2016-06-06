package handler

import (
	"html/template"
	"log"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {

	path := "views/content/404.tmpl"
	t, err := StandardTemplate(path)
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

func StandardTemplate(filepaths ...string) (*template.Template, error) {
	paths := append(filepaths, "templates/footer.tmpl",
		"templates/metatags.tmpl", "templates/logo.tmpl")
	t, err := template.ParseFiles(paths...)
	return t, err
}
