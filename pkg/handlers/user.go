package handlers

// TODO need to trim and verify usernames, passwords and emails.

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/authentication"
	"github.com/sprioc/sprioc-core/pkg/contentStorage"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"
)

type signUpFields struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) Response {
	decoder := json.NewDecoder(r.Body)

	newUser := signUpFields{}

	err := decoder.Decode(&newUser)
	if err != nil {
		return Resp("Bad Request", http.StatusBadRequest)
	}

	if mongo.ExistsUserName(newUser.Username) || mongo.ExistsEmail(newUser.Email) {
		return Resp("Username or Email already exist", http.StatusConflict)
	}

	password, salt, err := authentication.GetSaltPass(newUser.Password)
	if err != nil {
		return Resp("Error adding user", http.StatusConflict)
	}

	usr := model.User{
		ID:        bson.NewObjectId(),
		Email:     newUser.Email,
		Pass:      password,
		Salt:      salt,
		ShortCode: newUser.Username,
	}

	err = store.CreateUser(mongo, usr)
	if err != nil {
		return Resp("Error adding user", http.StatusConflict)
	}

	return Response{Code: http.StatusAccepted}
}

func AvatarUpload(w http.ResponseWriter, r *http.Request) Response {
	user, userRef, err := getLoggedInUser(r)
	if err != nil {
		return err.(Response)
	}

	file, err := ioutil.ReadAll(r.Body)
	n := len(file)

	if n == 0 {
		return Resp("Cannot upload file with 0 bytes.", http.StatusBadRequest)
	}

	err = contentStorage.ProccessImage(file, n, user.ShortCode, "avatar")
	if err != nil {
		log.Println(err)
		return Resp(err.Error(), http.StatusBadRequest)
	}

	sources := formatAvatarSources(user.ShortCode)

	err = store.ModifyAvatar(mongo, userRef, sources)
	if err != nil {
		return Resp("Unable to add image", http.StatusInternalServerError)
	}
	return Response{Code: http.StatusAccepted}
}

func formatAvatarSources(shortcode string) model.ImgSource {
	const prefix = "https://images.sprioc.xyz/avatars/"
	var resourceBaseURL = prefix + shortcode
	return model.ImgSource{
		Raw:    resourceBaseURL,
		Large:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
		Medium: resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
		Small:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
		Thumb:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) Response {
	UID := mux.Vars(r)["username"]

	user, err := store.GetUser(mongo, GetUserRef(UID))
	if err != nil {
		return Resp("Not Found", http.StatusNotFound)
	}

	dat, err := json.Marshal(user)
	if err != nil {
		return Resp("Unable to write JSON", http.StatusInternalServerError)
	}

	return Response{Code: http.StatusOK, Data: dat}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) Response {
	return executeCommand(w, r, getUserInterface, store.DeleteUser)
}

func FavoriteUser(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getUserInterface, store.FavoriteUser)
}

func FollowUser(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getUserInterface, store.FollowUser)
}

func UnFavoriteUser(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getUserInterface, store.UnFavoriteUser)
}

func UnFollowUser(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getUserInterface, store.UnFollowUser)
}

func ModifyUser(w http.ResponseWriter, r *http.Request) Response {
	username := mux.Vars(r)["username"]
	ref := GetUserRef(username)
	return executeCheckedModification(r, ref)
}
