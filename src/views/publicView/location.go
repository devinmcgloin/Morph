package publicView

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func LocationView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	LID, err := strconv.Atoi(ps.ByName("LID"))
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}

	locations, err := SQL.GetCollectionViewByLocation(uint64(LID))
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	t, err := common.StandardTemplate("templates/pages/locationView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, locations)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

}
