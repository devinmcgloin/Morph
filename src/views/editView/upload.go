package editView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/api/auth"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
)

func UploadView(w http.ResponseWriter, r *http.Request) {
	loggedIn, user := auth.CheckUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", 301)
		return
	}

	t, err := common.StandardTemplate("templates/pages/editView.tmpl")
	if err != nil {
		log.Printf("Error while parsing template: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	defaultView := model.DefaultView{
		Auth: user,
	}

	err = t.Execute(w, defaultView)
	if err != nil {
		log.Printf("Error while executing template %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
