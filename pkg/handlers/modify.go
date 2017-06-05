package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
)

func AddImageTag(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]
	tag := strings.TrimSpace(vars["tag"])

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

	tagRef, resp := core.GetTagRef(tag)
	if !resp.Ok() {
		return resp
	}

	return core.AddImageTag(usrRef, imageRef, tagRef)
}

func RemoveImageTag(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]
	tag := strings.TrimSpace(vars["tag"])

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

	tagRef, resp := core.GetTagRef(tag)
	if !resp.Ok() {
		return resp
	}

	return core.RemoveImageTag(usrRef, imageRef, tagRef)
}

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
