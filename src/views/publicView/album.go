package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
)

func AlbumView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["username"]
	albumTitle := vars["Album"]

	album, err := mongo.GetAlbumCollectionView(userName, albumTitle)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	t, err := common.StandardTemplate("templates/pages/albumView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, album)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

}
