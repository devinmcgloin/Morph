package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/morphError"
	"github.com/devinmcgloin/morph/src/views/common"
)

func MostRecentView(w http.ResponseWriter, r *http.Request) error {

	var images model.CollectionView

	images, err := mongo.GetNumMostRecentImg(10)
	if err != nil {
		return morphError.New(err, "Unable to get collection", 523)

	}

	usr, valid := session.GetUser(r)
	if valid {
		images.Auth = usr
	}
	return common.ExecuteTemplate(w, r, "templates/public/index.tmpl", images)
}
