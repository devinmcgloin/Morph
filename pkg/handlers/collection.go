package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/store"
)

func GetCollection(w http.ResponseWriter, r *http.Request) Response {
	CID := mux.Vars(r)["CID"]

	user, err := store.GetCollection(mongo, GetCollectionRef(CID))
	if err != nil {
		return Resp("Not Found", http.StatusNotFound)
	}

	dat, err := json.Marshal(user)
	if err != nil {
		return Resp("Unable to write JSON", http.StatusInternalServerError)
	}

	return Response{Code: http.StatusOK, Data: dat}
}

func CreateCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented}
}

// TODO these need to pull targets from request body
func AddImageToCollection(w http.ResponseWriter, r *http.Request) Response {
	CID := mux.Vars(r)["CID"]
	cRef := GetCollectionRef(CID)
	IID := mux.Vars(r)["IID"]
	iRef := GetCollectionRef(IID)
	return executeUncheckedModification(r, bson.M{"$push": bson.M{"images": iRef}}, cRef)
}

func AddUserToCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented}

}

func FavoriteCollection(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getCollectionInterface, store.FavoriteCollection)
}

func FollowCollection(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getCollectionInterface, store.FollowCollection)
}

func DeleteCollection(w http.ResponseWriter, r *http.Request) Response {
	return executeCommand(w, r, getCollectionInterface, store.DeleteCollection)
}

// TODO these need to pull targets from request body
func DeleteImageFromCollection(w http.ResponseWriter, r *http.Request) Response {
	CID := mux.Vars(r)["CID"]
	cRef := GetCollectionRef(CID)
	IID := mux.Vars(r)["IID"]
	iRef := GetCollectionRef(IID)
	return executeUncheckedModification(r, bson.M{"$push": bson.M{"images": iRef}}, cRef)
}

func DeleteUserFromCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented}
}

func UnFavoriteCollection(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getCollectionInterface, store.UnFavoriteAlbum)
}

func UnFollowCollection(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getCollectionInterface, store.UnFollowCollection)
}
