package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
)

func LocationView(w http.ResponseWriter, r *http.Request) error {

	var locations model.LocCollectionView

	usr, valid := session.GetUser(r)
	if valid {
		locations.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/public/UnderConstruction.tmpl", locations)
}
