package core

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/conductor/pkg/model"
	"github.com/sprioc/conductor/pkg/refs"
	"github.com/sprioc/conductor/pkg/rsp"
	"github.com/sprioc/conductor/pkg/store"
)

// TODO images needs to have access to the collection.
func AddImageToCollection(requestFrom model.User, col model.DBRef, additions map[string][]string) rsp.Response {
	if col.Collection != "collections" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	var links []string
	var ok bool

	if links, ok = additions["images"]; !ok {
		return rsp.Response{Message: "Invalid body", Code: http.StatusBadRequest}
	}

	if !inRef(col, requestFrom.Collections) {
		return rsp.Response{Message: "User cannot delete collection they do not own.", Code: http.StatusUnauthorized}
	}

	if !store.ExistsCollectionID(col.Shortcode) {
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

func DeleteImageFromCollection(requestFrom model.User, col model.DBRef, deletions map[string][]string) rsp.Response {
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

func AddTagsToImage(requestFrom model.User, ref model.DBRef, additions map[string][]string) rsp.Response {
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

func RemoveTagsFromImage(requestFrom model.User, ref model.DBRef, deletions map[string][]string) rsp.Response {
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
