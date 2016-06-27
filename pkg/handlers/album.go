package handlers

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/store"
)

func GetAlbum(w http.ResponseWriter, r *http.Request) Response {
	CID := mux.Vars(r)["CID"]

	album, err := store.GetAlbum(mongo, GetAlbumRef(CID))
	if err != nil {
		return Resp("Not Found", http.StatusNotFound)
	}

	return Response{Code: http.StatusOK, Data: album}
}

func CreateAlbum(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented}
}

// TODO these need to pull targets from request body
func AddImageToAlbum(w http.ResponseWriter, r *http.Request) Response {
	AID := mux.Vars(r)["AID"]
	aRef := GetAlbumRef(AID)
	IID := mux.Vars(r)["IID"]
	iRef := GetAlbumRef(IID)
	return executeUncheckedModification(r, bson.M{"$push": bson.M{"images": iRef}}, aRef)
}

func FavoriteAlbum(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getAlbumInterface, store.FavoriteAlbum)
}

func FollowAlbum(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getAlbumInterface, store.FollowAlbum)
}

func DeleteAlbum(w http.ResponseWriter, r *http.Request) Response {
	return executeCommand(w, r, getAlbumInterface, store.DeleteAlbum)
}

// TODO these need to pull targets from request body
func DeleteImageFromAlbum(w http.ResponseWriter, r *http.Request) Response {
	AID := mux.Vars(r)["AID"]
	aRef := GetAlbumRef(AID)
	IID := mux.Vars(r)["IID"]
	iRef := GetAlbumRef(IID)
	return executeUncheckedModification(r, bson.M{"$pull": bson.M{"images": iRef}}, aRef)
}

func UnFavoriteAlbum(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getAlbumInterface, store.UnFavoriteAlbum)
}

func UnFollowAlbum(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getAlbumInterface, store.UnFollowAlbum)
}
