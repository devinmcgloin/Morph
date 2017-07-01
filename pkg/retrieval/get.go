package retrieval

import (
	"net/http"

	"strconv"

	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/sql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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

func GetUserRef(username string) (model.Ref, rsp.Response) {
	usr, err := sql.GetUserRef(username)
	if err != nil {
		return model.Ref{}, rsp.Response{Code: http.StatusNotFound, Message: "Unable to retrieve reference"}
	}

	return usr, rsp.Response{Code: http.StatusOK}
}

func GetImageRef(shortcode string) (model.Ref, rsp.Response) {
	usr, err := sql.GetImageRef(shortcode)
	if err != nil {
		return model.Ref{}, rsp.Response{Code: http.StatusNotFound, Message: "Unable to retrieve reference"}
	}

	return usr, rsp.Response{Code: http.StatusOK}
}

func GetUser(ref model.Ref) (model.User, rsp.Response) {
	if ref.Collection != model.Users {
		return model.User{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	user, err := sql.GetUser(ref.Id, true)
	if err != nil {
		switch err.Error() {
		case "User not found.":
			return model.User{}, rsp.Response{Message: err.Error(), Code: http.StatusNotFound}
		default:
			return model.User{}, rsp.Response{Message: err.Error(), Code: http.StatusInternalServerError}
		}
	}
	return user, rsp.Response{Code: http.StatusOK}
}

func GetImage(ref model.Ref) (model.Image, rsp.Response) {
	if ref.Collection != model.Images {
		return model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	image, err := sql.GetImage(ref.Id)
	if err != nil {
		return model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return image, rsp.Response{Code: http.StatusOK}
}

func GetUserFollowed(ref model.Ref) ([]model.User, rsp.Response) {
	if ref.Collection != model.Users {
		return []model.User{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}
	images, err := sql.GetUserFollowed(ref.Id)
	if err != nil {
		return []model.User{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}

func GetUserFavorites(ref model.Ref) ([]model.Image, rsp.Response) {
	if ref.Collection != model.Users {
		return []model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}
	images, err := sql.GetUserFavorites(ref.Id)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}

func GetUserImages(ref model.Ref) ([]model.Image, rsp.Response) {
	if ref.Collection != model.Users {
		return []model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}
	images, err := sql.GetUserImages(ref.Id)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}

func GetRecentImages(limit int) ([]model.Image, rsp.Response) {
	images, err := sql.GetRecentImages(limit)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}

func GetFeaturedImages(limit int) ([]model.Image, rsp.Response) {
	images, err := sql.GetFeaturedImages(limit)
	if err != nil {
		return []model.Image{}, rsp.Response{Message: err.Error(),
			Code: http.StatusInternalServerError}
	}
	return images, rsp.Response{Code: http.StatusOK}
}
