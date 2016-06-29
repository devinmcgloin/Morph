package core

import (
	"net/http"
	"reflect"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/refs"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"github.com/sprioc/sprioc-core/pkg/store"
)

// REVIEW this code is narly and slow, not sure how well the db will maintain
// internal consistency. Mongo doesn't seem to give many gurantees.

func DeleteImage(requestFrom model.User, imageRef model.DBRef) rsp.Response {
	if imageRef.Collection != "images" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	if !inRef(imageRef, requestFrom.Images) {
		return rsp.Response{Message: "User cannot delete image they do not own.", Code: http.StatusUnauthorized}
	}

	image, resp := GetImage(imageRef)
	if !resp.Ok() {
		return rsp.Response{Code: http.StatusNotFound}
	}

	for _, usr := range image.FavoritedBy {
		UnFavorite(usr, imageRef)
	}

	for _, col := range image.Collections {
		Modify(col, bson.M{"$pull": bson.M{"collections": imageRef}})
	}

	resp = Modify(refs.GetUserRef(requestFrom.ShortCode), bson.M{"$pull": bson.M{"images": imageRef}})
	if !resp.Ok() {
		return rsp.Response{Message: "Unable to delete image", Code: http.StatusInternalServerError}
	}

	err := store.Delete(imageRef)
	if err != nil {
		return rsp.Response{Message: "Unable to delete image", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}

}

func DeleteCollection(requestFrom model.User, collectionRef model.DBRef) rsp.Response {
	if collectionRef.Collection != "collections" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	if !inRef(collectionRef, requestFrom.Collections) {
		return rsp.Response{Message: "User cannot delete collection they do not own.", Code: http.StatusUnauthorized}
	}

	collection, resp := GetCollection(collectionRef)
	if !resp.Ok() {
		return rsp.Response{Code: http.StatusNotFound}
	}

	for _, usr := range collection.FavoritedBy {
		UnFavorite(usr, collectionRef)
	}

	for _, usr := range collection.FollowedBy {
		UnFavorite(usr, collectionRef)
	}

	for _, img := range collection.Images {
		Modify(img, bson.M{"$pull": bson.M{"collections": collectionRef}})
	}

	resp = Modify(refs.GetUserRef(requestFrom.ShortCode), bson.M{"$pull": bson.M{"collections": collectionRef}})
	if !resp.Ok() {
		return rsp.Response{Message: "Unable to delete collection", Code: http.StatusInternalServerError}
	}

	err := store.Delete(collectionRef)
	if err != nil {
		return rsp.Response{Message: "Unable to delete collection", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func DeleteUser(requestFrom model.User, user model.DBRef) rsp.Response {
	if user.Collection != "users" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	if requestFrom.ShortCode != user.Shortcode {
		return rsp.Response{Message: "User cannot users besides themselves", Code: http.StatusUnauthorized}
	}

	for _, image := range requestFrom.Images {
		DeleteImage(requestFrom, image)
	}

	for _, col := range requestFrom.Collections {
		DeleteCollection(requestFrom, col)
	}

	for _, follow := range requestFrom.Followes {
		UnFollow(user, follow)
	}

	for _, followerRef := range requestFrom.FollowedBy {
		UnFollow(followerRef, user)
	}

	for _, favorite := range requestFrom.Favorites {
		UnFavorite(refs.GetUserRef(requestFrom.ShortCode), favorite)
	}

	for _, favoritedBy := range requestFrom.FavoritedBy {
		UnFavorite(favoritedBy, user)
	}

	err := store.Delete(user)
	if err != nil {
		return rsp.Response{Message: "Cannot delete user", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func inRef(item model.DBRef, collection []model.DBRef) bool {
	for _, x := range collection {
		if reflect.DeepEqual(x, item) {
			return true
		}
	}
	return false
}
