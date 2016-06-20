package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
)

func AlbumView(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userName := vars["username"]
	albumTitle := vars["Album"]

	album, err := mongo.GetAlbumCollectionView(userName, albumTitle)
	if err != nil {
		return morphError.New(err, "Unable to get collection", 523)
	}

	usr, valid := session.GetUser(r)
	if valid {
		album.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/public/albumView.tmpl", album)
}
