package modification

import (
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/sql"
)

func AddImageTag(usr model.Ref, image model.Ref, tag model.Ref) rsp.Response {
	resp := permission(usr, model.CanEdit, image)
	if !resp.Ok() {
		return resp
	}

	err := sql.AddTag(image.Id, tag.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to modify image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusCreated}
}

func RemoveImageTag(usr model.Ref, image model.Ref, tag model.Ref) rsp.Response {
	resp := permission(usr, model.CanEdit, image)
	if !resp.Ok() {
		return resp
	}

	err := sql.RemoveTag(image.Id, tag.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to modify image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusCreated}
}

func FeatureImage(usr model.Ref, image model.Ref) rsp.Response {
	resp := permission(usr, model.CanEdit, image)
	if !resp.Ok() {
		return resp
	}

	err := sql.Feature(image.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to feature image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func UnFeatureImage(usr model.Ref, image model.Ref) rsp.Response {
	resp := permission(usr, model.CanEdit, image)
	if !resp.Ok() {
		return resp
	}

	err := sql.UnFeature(image.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to feature image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func FavoriteImage(usr model.Ref, image model.Ref) rsp.Response {
	resp := permission(usr, model.CanView, image)
	if !resp.Ok() {
		return resp
	}

	err := sql.Favorite(usr.Id, image.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to feature image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func UnFavoriteImage(usr model.Ref, image model.Ref) rsp.Response {
	resp := permission(usr, model.CanView, image)
	if !resp.Ok() {
		return resp
	}

	err := sql.UnFavorite(usr.Id, image.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to feature image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}
func FollowUser(usr model.Ref, userB model.Ref) rsp.Response {
	resp := permission(usr, model.CanView, userB)
	if !resp.Ok() {
		return resp
	}

	err := sql.Follow(usr.Id, userB.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to feature image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func UnFollowUser(usr model.Ref, userB model.Ref) rsp.Response {
	resp := permission(usr, model.CanView, userB)
	if !resp.Ok() {
		return resp
	}

	err := sql.UnFollow(usr.Id, userB.Id)
	if err != nil {
		return rsp.Response{Message: "Unable to feature image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}
