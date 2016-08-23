package handlers

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

func UnFavoriteCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to unfavorite a collection"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	return core.UnFavorite(userRef, ref)
}

func FavoriteCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to favorite a collection"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	return core.Favorite(userRef, ref)
}

func UnFavoriteUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to unfavorite a user"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	return core.UnFavorite(userRef, ref)
}

func FavoriteUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to favorite a user"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	return core.Favorite(userRef, ref)
}

func UnFavoriteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to unfavorite a image"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	return core.UnFavorite(userRef, ref)
}

func FavoriteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to favorite a image"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	return core.Favorite(userRef, ref)
}
