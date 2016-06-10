package publicView

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devinmcgloin/morph/src/api/SQL"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func AlbumView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	AID, err := strconv.Atoi(ps.ByName("AID"))
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}
	log.Printf("Accessing user:%d", uint64(AID))
	album, err := SQL.GetAlbumCollectionView(uint64(AID))
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
