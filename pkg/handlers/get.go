package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
)

func GetUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["username"]

	ref, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	user, resp := core.GetUser(ref)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusOK, Data: user}
}

func GetLoggedInUser(w http.ResponseWriter, r *http.Request) rsp.Response {

	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to use this endpoint"}
	}

	usrRef := val.(model.Ref)
	user, resp := core.GetUser(usrRef)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusOK, Data: user}
}

func GetImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["IID"]

	ref, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	img, resp := core.GetImage(ref)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusOK, Data: img}
}
func GetUserFollowed(w http.ResponseWriter, r *http.Request) rsp.Response {
	vars := mux.Vars(r)
	username := vars["username"]

	ref, resp := core.GetUserRef(username)
	if !resp.Ok() {
		return resp
	}

	users, resp := core.GetUserFollowed(ref)
	if !resp.Ok() {
		return resp
	}
	return rsp.Response{Code: http.StatusOK, Data: users}
}
func GetUserFavorites(w http.ResponseWriter, r *http.Request) rsp.Response {
	vars := mux.Vars(r)
	username := vars["username"]

	ref, resp := core.GetUserRef(username)
	if !resp.Ok() {
		return resp
	}

	imgs, resp := core.GetUserFavorites(ref)
	if !resp.Ok() {
		return resp
	}
	return rsp.Response{Code: http.StatusOK, Data: imgs}
}

func GetUserImages(w http.ResponseWriter, r *http.Request) rsp.Response {
	vars := mux.Vars(r)
	username := vars["username"]

	ref, resp := core.GetUserRef(username)
	if !resp.Ok() {
		return resp
	}

	imgs, resp := core.GetUserImages(ref)
	if !resp.Ok() {
		return resp
	}
	return rsp.Response{Code: http.StatusOK, Data: imgs}
}

func GetRecentImages(w http.ResponseWriter, r *http.Request) rsp.Response {
	vars := mux.Vars(r)
	// save to ignore error as route has to match [0-9]+ regex to hit his handler
	limit, _ := strconv.Atoi(vars["limit"])
	imgs, resp := core.GetRecentImages(limit)
	if !resp.Ok() {
		return resp
	}
	return rsp.Response{Code: http.StatusOK, Data: imgs}
}

func GetFeaturedImages(w http.ResponseWriter, r *http.Request) rsp.Response {
	vars := mux.Vars(r)
	// save to ignore error as route has to match [0-9]+ regex to hit his handler
	limit, _ := strconv.Atoi(vars["limit"])
	imgs, resp := core.GetFeaturedImages(limit)
	if !resp.Ok() {
		return resp
	}
	return rsp.Response{Code: http.StatusOK, Data: imgs}
}
