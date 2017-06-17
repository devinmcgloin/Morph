package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
)

func FeatureImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	imageRef, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.FeatureImage(usrRef, imageRef)
}

func UnFeatureImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	imageRef, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.UnFeatureImage(usrRef, imageRef)
}

func FavoriteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	imageRef, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.FavoriteImage(usrRef, imageRef)
}

func UnFavoriteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	imageRef, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.UnFavoriteImage(usrRef, imageRef)
}

func Follow(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["username"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	userRef, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.FollowUser(usrRef, userRef)
}

func UnFollow(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["username"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	userRef, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.UnFollowUser(usrRef, userRef)
}

func PatchImage(w http.ResponseWriter, r *http.Request) rsp.Response {

	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	image, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	decoder := json.NewDecoder(r.Body)

	var request map[string]interface{}

	err := decoder.Decode(&request)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.PatchImage(usrRef, image, request)
}

func PatchUser(w http.ResponseWriter, r *http.Request) rsp.Response {

	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["username"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	image, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	decoder := json.NewDecoder(r.Body)

	var request map[string]interface{}

	err := decoder.Decode(&request)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.PatchUser(usrRef, image, request)
}
