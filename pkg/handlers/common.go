package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/authorization"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"
	"gopkg.in/mgo.v2/bson"
)

//TODO need to split this into two difference handlers. One for internal use,
//and the other for external.

// TODO need to re eval the special functions in store for specific use cases.
// Not going to be used if things are routed through modifications.

// executeCheckedModification is for use on external api endpoints as it parses
// changes from the request body and checks that the changes are authorized.
func executeCheckedModification(r *http.Request, ref model.DBRef) Response {
	target, err := store.GetRef(mongo, ref)
	if err != nil {
		return Resp("Unable to retrieve document", http.StatusNotFound)
	}

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return Resp("Unauthorized Request", http.StatusUnauthorized)
	}

	var changes bson.M

	err = json.NewDecoder(r.Body).Decode(&changes)
	if err != nil {
		log.Println(err)
		return Resp("Malformed Request", http.StatusBadRequest)
	}

	err = authorization.VerifyChanges(user, target, changes, false)
	if err != nil {
		return Resp(err.Error(), http.StatusUnauthorized)
	}

	err = store.Modify(mongo, ref, changes)
	if err != nil {
		return Resp(err.Error(), http.StatusBadRequest)
	}

	return Response{Code: 200}
}

func executeUncheckedModification(r *http.Request, changes bson.M, ref model.DBRef) Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return Resp("Unauthorized Request", http.StatusUnauthorized)
	}

	doc, err := store.GetRef(mongo, ref)
	if err != nil {
		return Resp("Not Found", http.StatusNotFound)
	}

	err = authorization.VerifyChanges(user, doc, changes, true)
	if err != nil {
		return Resp(err.Error(), http.StatusUnauthorized)
	}

	err = store.Modify(mongo, ref, changes)
	if err != nil {
		return Resp(err.Error(), http.StatusBadRequest)
	}

	return Response{Code: 200}
}

func getImage(r *http.Request) (model.Image, model.DBRef, error) {
	id := mux.Vars(r)["IID"]

	ref := GetImageRef(id)

	img, err := store.GetImage(mongo, ref)
	if err != nil {
		return model.Image{}, model.DBRef{}, Resp("Image not found", http.StatusNotFound)
	}
	return img, ref, nil
}

func getLoggedInUser(r *http.Request) (model.User, model.DBRef, error) {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return model.User{}, model.DBRef{}, Resp("Unauthorized Request", http.StatusUnauthorized)
	}
	return user, GetUserRef(user.ShortCode), nil
}

func getUser(r *http.Request) (model.User, model.DBRef, error) {
	id := mux.Vars(r)["username"]

	ref := GetUserRef(id)

	usr, err := store.GetUser(mongo, ref)
	if err != nil {
		return model.User{}, model.DBRef{}, Resp("User not found", http.StatusNotFound)
	}
	return usr, ref, nil
}

func getAlbum(r *http.Request) (model.Album, model.DBRef, error) {
	id := mux.Vars(r)["AID"]

	ref := GetAlbumRef(id)

	alb, err := store.GetAlbum(mongo, ref)
	if err != nil {
		return model.Album{}, model.DBRef{}, Resp("Album not found", http.StatusNotFound)
	}
	return alb, ref, nil
}

func getCollection(r *http.Request) (model.Collection, model.DBRef, error) {
	id := mux.Vars(r)["CID"]

	ref := GetAlbumRef(id)

	col, err := store.GetCollection(mongo, ref)
	if err != nil {
		return model.Collection{}, model.DBRef{}, Resp("Album not found", http.StatusNotFound)
	}
	return col, ref, nil
}

// func executeCommand(w http.ResponseWriter, r *http.Request,
// 	userGetter func(r *http.Request) (model.Collection, model.DBRef, error),
// 	targetGetter func(r *http.Request) (model.Collection, model.DBRef, error),
// 	operation func(ds *MgoStore)(ref model.DBRef)error){
//
// }
