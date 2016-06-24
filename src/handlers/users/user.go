package users

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/sprioc/src/api/auth"
	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
)

var mongo = store.ConnectStore()

type signUpFields struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	newUser := signUpFields{}

	err := decoder.Decode(&newUser)
	if err != nil {
		return spriocError.New(err, "Bad Request", http.StatusBadRequest)
	}

	if mongo.ExistsUserName(newUser.Username) || mongo.ExistsEmail(newUser.Email) {
		return spriocError.New(err, "Username or Email already exist", http.StatusConflict)
	}

	password, salt, err := auth.GetSaltPass(newUser.Password)
	if err != nil {
		return spriocError.New(err, "Error adding user", http.StatusConflict)
	}

	id := mongo.GetNewUserID()

	usr := model.User{
		UserName: newUser.Username,
		Email:    newUser.Email,
		Pass:     password,
		Salt:     salt,
		ID:       id,
	}

	err = mongo.CreateUser(usr)
	if err != nil {
		return spriocError.New(err, "Error adding user", http.StatusConflict)
	}

	return nil
}

func AvatarUploadHander(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", http.StatusNotImplemented)
}

func formatSources(ID bson.ObjectId) model.ImgSource {
	const prefix = "https://images.sprioc.xyz/avatars/"
	var resourceBaseURL = prefix + ID.Hex()
	return model.ImgSource{
		Raw:    resourceBaseURL,
		Large:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
		Medium: resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
		Small:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
		Thumb:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) error {
	username := mux.Vars(r)["username"]

	user, err := mongo.GetByUserName(username)
	if err != nil {
		return spriocError.New(err, "Not Found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return spriocError.New(err, "Unable to write JSON", http.StatusInternalServerError)
	}

	return nil
}
