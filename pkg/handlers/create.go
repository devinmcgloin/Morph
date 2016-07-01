package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func CreateCollection(w http.ResponseWriter, r *http.Request) rsp.Response {

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to create a collection"}
	}

	user = val.(model.User)

	decoder := json.NewDecoder(r.Body)

	var newCollection map[string]string

	err := decoder.Decode(&newCollection)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.CreateCollection(user, newCollection)
}

func CreateUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	decoder := json.NewDecoder(r.Body)

	var newUser map[string]string

	err := decoder.Decode(&newUser)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.CreateUser(newUser)
}
