package auth

import (
	"log"
	"net/http"
	"strings"

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

	log.Println(gothic.GetState(r))
	gothic.GetProviderName = getProvider

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", user)

	http.Redirect(w, r, "/", 301)

}

func getProvider(r *http.Request) (string, error) {
	url := r.URL.String()
	provider := strings.Split(url, "/")[1]
	return provider, nil
}
