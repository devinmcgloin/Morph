package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/auth"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
)

func MostRecentView(w http.ResponseWriter, r *http.Request) {

	var images model.CollectionView

	loggedIn, user := auth.CheckUser(r)
	if loggedIn {
		images.Auth = user
	}

	images, err := mongo.GetNumMostRecentImg(10)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	common.ExecuteTemplate(w, r, "templates/pages/index.tmpl", images)
	return
}
