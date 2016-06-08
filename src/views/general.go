package views

import (
	"html/template"
	"io/ioutil"
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

	files, _ := ioutil.ReadDir("./templates")
	for _, f := range files {
		filepaths = append(filepaths, "templates/"+f.Name())
	}

	t, err := template.ParseFiles(filepaths...)
	return t, err
}
