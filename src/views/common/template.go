package common

import (
	"html/template"
	"io/ioutil"
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

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}) {

	t, err := StandardTemplate(template)
	if err != nil {
		SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		SomethingsWrong(w, r, err)
		return
	}
}
