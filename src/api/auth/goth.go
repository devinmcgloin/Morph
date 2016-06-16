package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

// BeginAuthHandler
func BeginAuthHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Entered BeginAuthHandler")
	gothic.BeginAuthHandler(w, r)
}

func UserLoginCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Entered UserLoginCallback")

	//TODO need to add the user account here and log them in.

	log.Println(gothic.GetState(r))
	gothic.GetProviderName = getProvider

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println(err)
		common.SomethingsWrong(w, r, err)
	}

	log.Printf("%v", user)

	http.Redirect(w, r, "/", 301)

}

func getProvider(r *http.Request) (string, error) {
	url := r.URL.String()
	provider := strings.Split(url, "/")[2]
	return provider, nil
}
