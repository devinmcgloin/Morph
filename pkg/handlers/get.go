package handlers

//
import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

//
// func GetCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
// 	id := mux.Vars(r)["CID"]
//
// 	ref := refs.GetCollectionRef(id)
//
// 	col, resp := core.GetCollection(ref)
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	user, resp := core.GetUser(col.Owner)
// 	if !resp.Ok() {
// 		return resp
// 	}
// 	log.Printf("%+v", col)
//
// 	refs.FillExternalCollection(&col, user)
//
// 	return rsp.Response{Code: http.StatusOK, Data: col}
// }
//
// func GetUser(w http.ResponseWriter, r *http.Request) rsp.Response {
// 	id := mux.Vars(r)["username"]
//
// 	ref := refs.GetUserRef(id)
//
// 	user, resp := core.GetUser(ref)
// 	if !resp.Ok() {
// 		return resp
// 	}
// 	refs.FillExternalUser(&user)
// 	return rsp.Response{Code: http.StatusOK, Data: user}
// }
//
// func GetLoggedInUser(w http.ResponseWriter, r *http.Request) rsp.Response {
//
// 	var user model.User
// 	val, ok := context.GetOk(r, "auth")
// 	if !ok {
// 		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to delete image"}
// 	}
//
// 	user = val.(model.User)
//
// 	refs.FillExternalUser(&user)
// 	return rsp.Response{Code: http.StatusOK, Data: user}
// }
//
func GetImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	img, resp := core.GetImage(ref)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusOK, Data: img}
}

//
// func IncrementDownloadCounter(w http.ResponseWriter, r *http.Request) rsp.Response {
// 	id := mux.Vars(r)["IID"]
//
// 	ref := refs.GetImageRef(id)
//
// 	resp := core.IncrementDownloads(ref)
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	return rsp.Response{Code: http.StatusOK}
//
// }
//
// func GetUserImages(w http.ResponseWriter, r *http.Request) rsp.Response {
// 	id := mux.Vars(r)["username"]
//
// 	ref := refs.GetUserRef(id)
//
// 	images, resp := core.GetUserImages(ref)
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	for _, img := range images {
// 		refs.FillExternalImage(img)
// 	}
//
// 	return rsp.Response{Code: http.StatusOK, Data: images}
// }
//
// func GetCollectionImages(w http.ResponseWriter, r *http.Request) rsp.Response {
//
// 	id := mux.Vars(r)["CID"]
//
// 	ref := refs.GetCollectionRef(id)
//
// 	images, resp := core.GetCollectionImages(ref)
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	if len(images) < 1 {
// 		return rsp.Response{Code: http.StatusOK, Data: images}
// 	}
//
// 	for _, img := range images {
// 		refs.FillExternalImage(img)
// 	}
//
// 	return rsp.Response{Code: http.StatusOK, Data: images}
//
// }
//
// func GetFeaturedImages(w http.ResponseWriter, r *http.Request) rsp.Response {
//
// 	images, resp := core.GetFeaturedImages()
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	for _, img := range images {
// 		refs.FillExternalImage(img)
// 	}
//
// 	return rsp.Response{Code: http.StatusOK, Data: images}
// }
