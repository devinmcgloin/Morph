package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"gopkg.in/mgo.v2/bson"
)

var mongo = store.NewStore()

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)

var s = securecookie.New(hashKey, blockKey)

var sessionLifetime = time.Minute * 10
var refreshAt = time.Minute * 1

// BeginAuthHandler
func BeginAuthHandler(w http.ResponseWriter, r *http.Request) error {
	gothic.BeginAuthHandler(w, r)
	return nil
}

func UserLoginCallback(w http.ResponseWriter, r *http.Request) error {

	gothic.GetProviderName = getProvider

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return spriocError.New(err, "Could not complete user Auth", 523)
	}

	internalUser := ConvertGothUser(user)
	err = RegisterUser(internalUser)
	if err != nil {
		return spriocError.New(err, "Could not register user", 523)

	}

	sessionID := session.NewSessionID()

	value := map[string]string{
		"provider":    internalUser.Provider,
		"provider_id": internalUser.ProviderID,
		"session_id":  sessionID,
	}

	if encoded, err := s.Encode("morph", value); err == nil {
		cookie := &http.Cookie{
			Name:    "morph",
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(sessionLifetime),
		}
		http.SetCookie(w, cookie)
	}

	internalUser, err = mongo.GetUserByUserName(internalUser.UserName)
	if err != nil {
		return spriocError.New(err, "Could not get user", 523)

	}

	err = session.SetUserID(sessionID, internalUser.ID)
	if err != nil {
		return spriocError.New(err, "Could not set user session", 523)
	}

	http.Redirect(w, r, "/", 302)
	return nil
}

func getProvider(r *http.Request) (string, error) {
	provider := mux.Vars(r)["provider"]
	return provider, nil
}

// CheckUser looks at the request, matches the cookie with the user and updates the
// cookie if it is close to expiration. Also returns the user object.
func CheckUser(r *http.Request) (bool, model.User) {

	sessionID, err := getSessionID(r)
	if err != nil {
		return false, model.User{}
	}

	userID, err := session.GetUserId(sessionID)
	if err != nil {
		return false, model.User{}
	}

	user, err := mongo.GetUserByID(userID)
	if err != nil {
		return false, model.User{}
	}
	return true, user
}

func RegisterUser(user model.User) error {
	userExists := mongo.ExistsUser(user.Provider, user.ProviderID)
	if userExists {
		return nil
	}
	user.ID = bson.NewObjectId()
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

func expireCookies(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	morphCookie, err := r.Cookie("morph")
	if err != nil {
		log.Println(err)
	} else {
		morphCookie.Expires = time.Now().Add(-time.Hour * 24)
		morphCookie.Value = "INVALID COOKIE"
		http.SetCookie(w, morphCookie)
	}

	gothicCookie, err := r.Cookie("_gothic_session")
	if err != nil {
		log.Println(err)
	} else {
		gothicCookie.Expires = time.Now().Add(-time.Hour * 24)
		gothicCookie.Value = "INVALID COOKIE"
		http.SetCookie(w, gothicCookie)
	}
	return w
}

func LogoutUser(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	sessionID, err := getSessionID(r)
	if err != nil {
		log.Println(err)
	}
	session.DeleteSessionID(sessionID)
	return expireCookies(w, r)
}

func RenewCookie(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	morphCookie, err := r.Cookie("morph")
	if err != nil {
		log.Println(err)

		// Checking if cookie has more than two minutes left
	} else if morphCookie.Expires.Before(time.Now().Add(refreshAt)) {
		return w
	} else {
		morphCookie.Expires = time.Now().Add(sessionLifetime)
		http.SetCookie(w, morphCookie)
	}
	return w
}

func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("morph")
	if err != nil {
		log.Println(err)
		return "", err
	}
	value := make(map[string]string)

	err = s.Decode("morph", cookie.Value, &value)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return value["session_id"], nil
}
