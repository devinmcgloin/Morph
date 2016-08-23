package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

func ModifyImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to feature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["IID"]
	ref := refs.GetImageRef(id)

	var changes bson.M
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&changes)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.ModifySecure(user, ref, changes)
}

func ModifyUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to feature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["username"]
	ref := refs.GetUserRef(id)

	var changes bson.M
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&changes)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.ModifySecure(user, ref, changes)
}

func ModifyCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to feature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["username"]
	ref := refs.GetCollectionRef(id)

	var changes bson.M
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&changes)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.ModifySecure(user, ref, changes)
}

func UnFeatureImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to unfeature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	return core.UnFeatureImage(user, ref)
}

func FeatureImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to feature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	return core.FeatureImage(user, ref)
}

func AddImageToCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to feature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	decoder := json.NewDecoder(r.Body)
	var additions map[string][]string

	err := decoder.Decode(&additions)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.AddImageToCollection(user, ref, additions)
}

func DeleteImageFromCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to feature a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["CID"]

	ref := refs.GetCollectionRef(id)

	decoder := json.NewDecoder(r.Body)
	var additions map[string][]string

	err := decoder.Decode(&additions)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.DeleteImageFromCollection(user, ref, additions)
}

func AddTagsToImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to add a tag to a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["CID"]

	ref := refs.GetImageRef(id)

	decoder := json.NewDecoder(r.Body)
	var additions map[string][]string

	err := decoder.Decode(&additions)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	resp := core.AddTagsToImage(user, ref, additions)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusAccepted}
}

func RemoveTagsFromImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to remove a tag from a image"}
	}

	user = val.(model.User)

	id := mux.Vars(r)["CID"]

	ref := refs.GetImageRef(id)

	decoder := json.NewDecoder(r.Body)
	var additions map[string][]string

	err := decoder.Decode(&additions)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	resp := core.RemoveTagsFromImage(user, ref, additions)
	if !resp.Ok() {
		return resp
	}

	return rsp.Response{Code: http.StatusAccepted}

}
