package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
)

func SearchView(w http.ResponseWriter, r *http.Request) error {

	var collection model.CollectionView

	usr, valid := session.GetUser(r)
	if valid {
		collection.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/public/UnderConstruction.tmpl", collection)

}
