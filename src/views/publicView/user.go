package publicView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func UserProfileView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	UserName := ps.ByName("UserName")

	log.Printf("Accessing user:%s", UserName)
	user, err := SQL.GetUserProfileView(UserName)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	log.Printf("%v", user.Images[0])

	t, err := common.StandardTemplate("templates/pages/userView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, user)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}
}
