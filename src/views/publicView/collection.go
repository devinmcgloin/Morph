package publicView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/api/auth"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func MostRecentView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	loggedIn, user := auth.CheckUser(r)
	log.Println(loggedIn, user)

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
