package publicView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
)

func UserProfileView(w http.ResponseWriter, r *http.Request) {

	UserName := mux.Vars(r)["username"]

	log.Printf("Accessing user:%s", UserName)
	user, err := mongo.GetUserProfileView(UserName)
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
