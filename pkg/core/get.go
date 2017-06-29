package core

import (
	"net/http"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
	"github.com/sprioc/composer/pkg/sql"
)

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
