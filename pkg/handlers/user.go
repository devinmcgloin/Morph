package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/sprioc/pkg/api/auth"
	"github.com/devinmcgloin/sprioc/pkg/model"
	"github.com/gorilla/mux"
)

type signUpFields struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) Response {
	decoder := json.NewDecoder(r.Body)

	newUser := signUpFields{}

	err := decoder.Decode(&newUser)
	if err != nil {
		return Resp("Bad Request", http.StatusBadRequest)
	}

	if mongo.ExistsUserName(newUser.Username) || mongo.ExistsEmail(newUser.Email) {
		return Resp("Username or Email already exist", http.StatusConflict)
	}

	password, salt, err := auth.GetSaltPass(newUser.Password)
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

	err = mongo.CreateUser(usr)
	if err != nil {
		return Resp("Error adding user", http.StatusConflict)
	}

	return Response{Code: http.StatusAccepted}
}

func AvatarUploadHander(w http.ResponseWriter, r *http.Request) Response {
	return Resp("Not Implemented", http.StatusNotImplemented)
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) Response {
	username := mux.Vars(r)["username"]

	user, err := mongo.GetByUserName(username)
	if err != nil {
		return Resp("Not Found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	dat, err := json.Marshal(user)
	if err != nil {
		return Resp("Unable to write JSON", http.StatusInternalServerError)
	}

	return Response{Code: http.StatusOK, Data: dat}
}
