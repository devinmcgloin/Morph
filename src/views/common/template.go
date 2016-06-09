package common

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func StandardTemplate(filepaths ...string) (*template.Template, error) {

	files, _ := ioutil.ReadDir("./templates/components/")
	for _, f := range files {
		filepaths = append(filepaths, "./templates/components/"+f.Name())
	}

	t, err := template.ParseFiles(filepaths...)
	return t, err
}

func RenderStatic(w http.ResponseWriter, r *http.Request, resources string) {
	t, err := StandardTemplate(resources)
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
