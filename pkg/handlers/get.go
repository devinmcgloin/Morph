package handlers

import (
	"log"
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

	user, resp := core.GetUser(col.Owner)
	if !resp.Ok() {
		return resp
	}
	log.Printf("%+v", col)

	refs.FillExternalCollection(&col, user)

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

func GetUserImages(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	images, resp := core.GetUserImages(ref)
	if !resp.Ok() {
		return resp
	}

	user, resp := core.GetUser(ref)
	if !resp.Ok() {
		return resp
	}

	for _, img := range images {
		refs.FillExternalImage(img, user)
	}

	return rsp.Response{Code: http.StatusOK, Data: images}
}

func GetCollectionImages(w http.ResponseWriter, r *http.Request) rsp.Response {

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	images, resp := core.GetCollectionImages(ref)
	if !resp.Ok() {
		return resp
	}

	log.Printf("%+v", images)

	if len(images) < 1 {
		return rsp.Response{Code: http.StatusOK, Data: images}
	}

	user, resp := core.GetUser(images[0].Owner)
	if !resp.Ok() {
		return resp
	}

	for _, img := range images {
		refs.FillExternalImage(img, user)
	}

	return rsp.Response{Code: http.StatusOK, Data: images}

}
