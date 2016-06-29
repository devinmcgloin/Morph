package core

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/refs"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"github.com/sprioc/sprioc-core/pkg/store"
)

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

	for _, link := range links {
		ref := refs.GetRef(link)
		err := store.Modify(col, bson.M{"$addToSet": bson.M{"images": ref}})
		if err != nil {
			return rsp.Response{Code: http.StatusInternalServerError}
		}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func DeleteImageToCollection(requestFrom model.User, ref model.DBRef, additons map[string][]string) rsp.Response {
	if ref.Collection != "collections" {
		return rsp.Response{Message: "Invalid reference", Code: http.StatusBadRequest}
	}

	if !inRef(ref, requestFrom.Collections) {
		return rsp.Response{Message: "User cannot delete collection they do not own.", Code: http.StatusUnauthorized}
	}

	if !store.ExistsCollectionID(ref.Shortcode) {
		return rsp.Response{Code: http.StatusNotFound}
	}

	err := store.Modify(ref, bson.M{"$pull": bson.M{"images": ref}})
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}
