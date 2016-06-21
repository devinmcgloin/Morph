package search

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
)

func CollectionTagView(w http.ResponseWriter, r *http.Request) error {

	tag := mux.Vars(r)["tag"]

	taggedImages, err := mongo.GetCollectionViewByTags(tag)
	if err != nil {
		return spriocError.New(err, "Unable to get collection", 523)

	}

	if len(taggedImages.Images) == 0 {
		return spriocError.New(err, "Collection was Empty", 404)

	}

	usr, valid := session.GetUser(r)
	if valid {
		taggedImages.Auth = usr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(taggedImages)

	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}
	return nil

}

func CollectionTagFeatureView(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", 404)
}
