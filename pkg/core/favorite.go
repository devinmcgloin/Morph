package core

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"gopkg.in/mgo.v2/bson"
)

// TODO would like to tell users the request failed due to the target not existsing.

func Favorite(user model.DBRef, target model.DBRef) rsp.Response {

	if user.Collection != "users" {
		return rsp.Response{Message: "Only users can favorite things.", Code: http.StatusBadRequest}
	}

	if user.Shortcode == target.Shortcode {
		return rsp.Response{Message: "User cannot favorite themselves.", Code: http.StatusBadRequest}
	}

	resp := Modify(target, bson.M{"$addToSet": bson.M{"favorited_by": user}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while adding relation", Code: http.StatusInternalServerError}
	}

	resp = Modify(user, bson.M{"$addToSet": bson.M{"favorites": target}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while adding relation", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}

}

func UnFavorite(user model.DBRef, target model.DBRef) rsp.Response {

	if user.Collection != "users" {
		return rsp.Response{Message: "Only users can unfavorite things.", Code: http.StatusBadRequest}
	}

	if user.Shortcode == target.Shortcode {
		return rsp.Response{Message: "User cannot unfavorite themselves.", Code: http.StatusBadRequest}
	}

	resp := Modify(target, bson.M{"$pull": bson.M{"favorited_by": user}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while removing relation", Code: http.StatusInternalServerError}
	}

	resp = Modify(user, bson.M{"$pull": bson.M{"favorites": target}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while removing relation", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
