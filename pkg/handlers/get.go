package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/refs"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func GetCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	col, resp := core.GetCollection(ref)
	if !resp.Ok() {
		return resp
	}
	return rsp.Response{Code: http.StatusOK, Data: col}
}

func GetUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	user, resp := core.GetUser(ref)
	if !resp.Ok() {
		return resp
	}
	refs.FillExternalUser(&user)
	return rsp.Response{Code: http.StatusOK, Data: user}
}

func GetImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	img, resp := core.GetImage(ref)
	if !resp.Ok() {
		return resp
	}

	user, resp := core.GetUser(img.Owner)
	if !resp.Ok() {
		return resp
	}

	refs.FillExternalImage(&img, user)
	return rsp.Response{Code: http.StatusOK, Data: img}
}
