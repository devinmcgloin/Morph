package auth

import (
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/morph/src/api/store"
	"github.com/devinmcgloin/morph/src/model"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

var mongo = store.NewStore()

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)
var s = securecookie.New(hashKey, blockKey)

// BeginAuthHandler
func BeginAuthHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Entered BeginAuthHandler")
	gothic.BeginAuthHandler(w, r)
}

func UserLoginCallback(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Entered UserLoginCallback")

	//TODO need to add the user account here and log them in.

	gothic.GetProviderName = getProvider

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println(err)
		common.SomethingsWrong(w, r, err)
		return
	}

	internalUser := ConvertGothUser(user)
	err = RegisterUser(internalUser)
	if err != nil {
		http.Redirect(w, r, "/", 301)
	}

	value := map[string]string{
		"provider":    internalUser.Provider,
		"provider_id": internalUser.ProviderID,
		"username":    internalUser.UserName,
	}

	if encoded, err := s.Encode("morph", value); err == nil {
		cookie := &http.Cookie{
			Name:  "morph",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}

	http.Redirect(w, r, "/", 301)
}

func getProvider(r *http.Request) (string, error) {
	url := r.URL.String()
	provider := strings.Split(url, "/")[2]
	return provider, nil
}

// CheckUser looks at the request, matches the cookie with the user and updates the
// cookie if it is close to expiration. Also returns the user object.
func CheckUser(r *http.Request) (bool, model.User) {
	cookie, err := r.Cookie("morph")
	if err != nil {
		log.Println(err)
		return false, model.User{}
	}
	value := make(map[string]string)

	err = s.Decode("morph", cookie.Value, &value)
	if err != nil {
		log.Println(err)
		return false, model.User{}
	}

	username := value["username"]
	user, err := mongo.GetUserByUserName(username)
	if err != nil {
		log.Println(err)
		return false, model.User{}
	}

	return true, user
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
	modelUser.UserName = user.NickName // TODO need to check username here
	modelUser.Name = user.Name
	modelUser.Bio = user.Description
	modelUser.ID = bson.NewObjectId()

	return modelUser
}
