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

func UnFollowCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to unfollow a collection"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	return core.UnFollow(userRef, ref)
}

func FollowCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to follow a collection"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	return core.Follow(userRef, ref)
}
func UnFollowUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to unfollow a user"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	return core.UnFollow(userRef, ref)
}

func FollowUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to Follow a user"}
	}

	user = val.(model.User)
	userRef := refs.GetUserRef(user.ShortCode)

	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	return core.Follow(userRef, ref)
}
