package handlers

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
)

func GetUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	user, resp := core.GetUser(ref)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusOK, Data: user}
}

func GetLoggedInUser(w http.ResponseWriter, r *http.Request) rsp.Response {

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to delete image"}
	}

	user = val.(model.User)

	return rsp.Response{Code: http.StatusOK, Data: user}
}

func GetImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	img, resp := core.GetImage(ref)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusOK, Data: img}
}
