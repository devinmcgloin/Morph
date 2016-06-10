package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth/gothic"
)

// BeginAuthHandler
func BeginAuthHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gothic.BeginAuthHandler(w, r)
}

func UserLoginCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(gothic.GetState(r))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		panic(err)
	}

	log.Printf("%v", user)

	http.Redirect(w, r, "/", 301)

}
