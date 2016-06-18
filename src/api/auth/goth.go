package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/devinmcgloin/morph/src/api/store"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

var mongo = store.NewStore()

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
		return
	}

	log.Printf("%v", user)

	http.Redirect(w, r, "/", 301)
}

func getProvider(r *http.Request) (string, error) {
	url := r.URL.String()
	provider := strings.Split(url, "/")[2]
	return provider, nil
}

func LoggedIn(r *http.Request) (bool, model.User) {
	return false, model.User{}
}

func RegisterUser(user model.User) error {
	userExists := mongo.ExistsUser(user.Provider, user.ProviderID)
	if userExists {
		return nil
	}
	err := mongo.AddUser(user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ConvertGothUser switches the user from goths interpertation to the internal one.
func ConvertGothUser(user goth.User) model.User {
	var modelUser model.User

	modelUser.Email = user.Email
	modelUser.Provider = user.Provider
	modelUser.UserName = user.NickName
	modelUser.Name = user.Name
	modelUser.Bio = user.Description

	return modelUser
}
