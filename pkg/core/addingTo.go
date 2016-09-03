package core

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/mongo"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

func AddImageToCollection(requestFrom model.Ref, col model.Ref, additions map[string][]string) rsp.Response {
	if col.ItemType != model.Collections {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var links []string
	var ok bool

	if links, ok = additions["images"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	if redis.Permissions(requestFrom, redis.CanEdit) {
		return rsp.Response{Message: "User cannot modify collection.", Code: http.StatusUnauthorized}
	}

	if !redis.Exists(col) {
		return rsp.Response{Code: http.StatusNotFound}
	}

	refs := refs.GetRefs(links)
	err := store.Modify(col, bson.M{"$addToSet": bson.M{"images": bson.M{"$each": refs}}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	for _, ref := range refs {
		err := store.Modify(ref, bson.M{"$addToSet": bson.M{"collections": col}})
		if err != nil {
			return rsp.Response{Code: http.StatusInternalServerError}
		}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func DeleteImageFromCollection(requestFrom model.Ref, col model.Ref, deletions map[string][]string) rsp.Response {
	if col.Collection != "collections" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var links []string
	var ok bool

	if links, ok = deletions["images"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	if !inRef(col, requestFrom.Collections) {
		return rsp.Response{Message: "User cannot delete collection they do not own.", Code: http.StatusUnauthorized}
	}

	if !store.ExistsCollectionID(col.Shortcode) {
		return rsp.Response{Code: http.StatusNotFound}
	}

	refs := refs.GetRefs(links)
	err := store.Modify(col, bson.M{"$pull": bson.M{"images": bson.M{"$each": refs}}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	for _, ref := range refs {
		err := store.Modify(ref, bson.M{"$pull": bson.M{"collections": col}})
		if err != nil {
			return rsp.Response{Code: http.StatusInternalServerError}
		}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func AddTagsToImage(requestFrom model.User, ref model.Ref, additions map[string][]string) rsp.Response {
	if ref.Collection != "images" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var tags []string
	var ok bool

	if tags, ok = additions["tags"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	if !inRef(ref, requestFrom.Images) {
		return rsp.Response{Message: "User cannot delete image they do not own.", Code: http.StatusUnauthorized}
	}

	if !store.ExistsImageID(ref.Shortcode) {
		return rsp.Response{Code: http.StatusNotFound}
	}

	err := store.Modify(ref, bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tags}}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func RemoveTagsFromImage(requestFrom model.User, ref model.Ref, deletions map[string][]string) rsp.Response {
	if ref.Collection != "images" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var tags []string
	var ok bool

	if tags, ok = deletions["tags"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	if !inRef(ref, requestFrom.Images) {
		return rsp.Response{Message: "User cannot delete image they do not own.", Code: http.StatusUnauthorized}
	}

	if !store.ExistsImageID(ref.Shortcode) {
		return rsp.Response{Code: http.StatusNotFound}
	}

	err := store.Modify(ref, bson.M{"$pull": bson.M{"tags": bson.M{"$each": tags}}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted}
}
