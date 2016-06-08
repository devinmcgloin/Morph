package endpoint

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func LoggedIn(r *http.Request) bool {
	return true

}

func LoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.Redirect(w, r, "/morph/", 302)
}

func SetAuthBoss() {

}
