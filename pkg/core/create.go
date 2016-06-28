package core

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"github.com/sprioc/sprioc-core/pkg/store"
)

func CreateImage(img model.Image) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented}
}

func CreateUser(userData map[string]string) rsp.Response {
	var username, email, password string
	var ok bool

	if username, ok = userData["username"]; !ok {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if email, ok = userData["email"]; !ok {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if password, ok = userData["password"]; !ok {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if store.ExistsUserName(username) || store.ExistsEmail(email) {
		return rsp.Response{Message: "Username or Email already exist", Code: http.StatusConflict}
	}

	securePassword, salt, err := GetSaltPass(password)
	if err != nil {
		return rsp.Response{Message: "Error adding user", Code: http.StatusConflict}
	}

	usr := model.User{
		ID:        bson.NewObjectId(),
		Email:     email,
		Pass:      securePassword,
		Salt:      salt,
		ShortCode: username,
		AvatarURL: formatSources("default"),
	}

	err = store.Create("users", usr)
	if err != nil {
		return rsp.Response{Message: "Error adding user", Code: http.StatusConflict}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
