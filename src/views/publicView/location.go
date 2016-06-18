package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func LocationView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	locShortText := ps.ByName("LocShortText")

	locations, err := mongo.GetCollectionViewByLocations(locShortText)
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
