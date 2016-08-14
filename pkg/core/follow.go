package core

import (
	"net/http"

	"github.com/sprioc/conductor/pkg/model"
	"github.com/sprioc/conductor/pkg/rsp"
	"gopkg.in/mgo.v2/bson"
)

// TODO would like to tell users the request failed due to the target not existsing.

func Follow(user model.DBRef, target model.DBRef) rsp.Response {

	if user.Shortcode == "" || target.Shortcode == "" {
		return rsp.Response{Code: http.StatusBadRequest, Message: "Shortcodes must not be empty"}
	}

	if user.Collection != "users" {
		return rsp.Response{Message: "Only users can unfollow things.", Code: http.StatusBadRequest}
	}

	if target.Collection == "images" {
		return rsp.Response{Message: "Users cannot follow images", Code: http.StatusBadRequest}
	}

	if user.Shortcode == target.Shortcode {
		return rsp.Response{Message: "User cannot follow themselves.", Code: http.StatusBadRequest}
	}

	resp := Modify(target, bson.M{"$addToSet": bson.M{"followed_by": user}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while adding relation", Code: http.StatusInternalServerError}
	}

	resp = Modify(user, bson.M{"$addToSet": bson.M{"followes": target}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while adding relation", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func UnFollow(user model.DBRef, target model.DBRef) rsp.Response {

	if user.Collection != "users" {
		return rsp.Response{Message: "Only users can unfollow things.", Code: http.StatusBadRequest}
	}

	if target.Collection == "images" {
		return rsp.Response{Message: "Users cannot unfollow images", Code: http.StatusBadRequest}
	}

	if user.Shortcode == target.Shortcode {
		return rsp.Response{Message: "User cannot unfollow themselves.", Code: http.StatusBadRequest}
	}

	resp := Modify(target, bson.M{"$pull": bson.M{"followed_by": user}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while removing relation", Code: http.StatusInternalServerError}
	}

	resp = Modify(user, bson.M{"$pull": bson.M{"followes": target}})
	if !resp.Ok() {
		return rsp.Response{Message: "Error while removing relation", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
