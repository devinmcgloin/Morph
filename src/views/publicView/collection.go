package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func MostRecentView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	locations, err := mongo.GetNumMostRecentImg(10)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	t, err := common.StandardTemplate("templates/pages/index.tmpl")
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
