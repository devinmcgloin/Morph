package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/authorization"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"
	"gopkg.in/mgo.v2/bson"
)

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

func getImageInterface(r *http.Request) (interface{}, model.DBRef, error) {
	id := mux.Vars(r)["IID"]

	ref := GetImageRef(id)

	img, err := store.GetImage(mongo, ref)
	if err != nil {
		return model.Image{}, model.DBRef{}, Resp("Image not found", http.StatusNotFound)
	}
	return img, ref, nil
}

func getImage(r *http.Request) (model.Image, model.DBRef, error) {
	doc, docRef, err := getImageInterface(r)
	convertDoc, ok := doc.(model.Image)
	if ok {
		return convertDoc, docRef, err
	}
	return model.Image{}, model.DBRef{}, errors.New("Unable to convert image")
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

func getUserInterface(r *http.Request) (interface{}, model.DBRef, error) {
	id := mux.Vars(r)["username"]

	ref := GetUserRef(id)

	usr, err := store.GetUser(mongo, ref)
	if err != nil {
		return model.User{}, model.DBRef{}, Resp("User not found", http.StatusNotFound)
	}
	return usr, ref, nil
}

func getUser(r *http.Request) (model.User, model.DBRef, error) {
	doc, docRef, err := getUserInterface(r)
	convertDoc, ok := doc.(model.User)
	if ok {
		return convertDoc, docRef, err
	}
	return model.User{}, model.DBRef{}, errors.New("Unable to convert user")
}

func getCollectionInterface(r *http.Request) (interface{}, model.DBRef, error) {
	id := mux.Vars(r)["CID"]

	ref := GetAlbumRef(id)

	col, err := store.GetCollection(mongo, ref)
	if err != nil {
		return model.Collection{}, model.DBRef{}, Resp("Album not found", http.StatusNotFound)
	}
	return col, ref, nil
}

func getCollection(r *http.Request) (model.Collection, model.DBRef, error) {
	doc, docRef, err := getCollectionInterface(r)
	convertDoc, ok := doc.(model.Collection)
	if ok {
		return convertDoc, docRef, err
	}
	return model.Collection{}, model.DBRef{}, errors.New("Unable to convert collection")
}

type getter func(r *http.Request) (interface{}, model.DBRef, error)
type opt func(ds *store.MgoStore, ref model.DBRef) error
type biDirectOpt func(ds *store.MgoStore, actor model.DBRef, subject model.DBRef) error

func executeCommand(w http.ResponseWriter, r *http.Request,
	targetGetter getter,
	operation opt) Response {

	target, targetRef, err := targetGetter(r)
	if err != nil {
		return err.(Response)
	}

	user, _, err := getLoggedInUser(r)
	if err != nil {
		return err.(Response)
	}

	_, err = authorization.Authorized(user, target)
	if err != nil {
		return Resp(err.Error(), http.StatusUnauthorized)
	}

	err = operation(mongo, targetRef)
	if err != nil {
		return Resp("Internal Server Error", http.StatusInternalServerError)
	}
	return Response{Code: http.StatusAccepted}
}

func executeBiDirectCommand(w http.ResponseWriter, r *http.Request, targetGetter getter,
	operation biDirectOpt) Response {

	target, targetRef, err := targetGetter(r)
	if err != nil {
		return err.(Response)
	}

	user, userRef, err := getLoggedInUser(r)
	if err != nil {
		return err.(Response)
	}

	_, err = authorization.Authorized(user, target)
	if err != nil {
		return Resp(err.Error(), http.StatusUnauthorized)
	}

	err = operation(mongo, userRef, targetRef)
	if err != nil {
		return Resp("Internal Server Error", http.StatusInternalServerError)
	}
	return Response{Code: http.StatusAccepted}
}
