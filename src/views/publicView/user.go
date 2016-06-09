package publicView

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func UserProfileView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	UID, err := strconv.Atoi(ps.ByName("UID"))
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}
	log.Printf("Accessing user:%d", uint64(UID))
	img, err := SQL.GetUserProfileView(uint64(UID))
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	t, err := common.StandardTemplate("templates/pages/userView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, img)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}
}
