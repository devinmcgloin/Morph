package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func AlbumView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userName := ps.ByName("UserName")
	albumTitle := ps.ByName("Album")

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
