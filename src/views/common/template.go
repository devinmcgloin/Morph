package common

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/devinmcgloin/morph/src/morphError"
)

func StandardTemplate(filepaths ...string) (*template.Template, error) {

	files, _ := ioutil.ReadDir("./templates/components/")
	for _, f := range files {
		filepaths = append(filepaths, "./templates/components/"+f.Name())
	}

	t, err := template.ParseFiles(filepaths...)
	return t, err
}

func ExecuteTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}) error {
	t, err := StandardTemplate(template)
	if err != nil {
		return morphError.New(err, "Unable to parse template", 523)
	}

	err = t.Execute(w, data)
	if err != nil {
		return morphError.New(err, "Unable to execute template", 523)
	}
	return nil
}
