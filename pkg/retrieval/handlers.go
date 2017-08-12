package retrieval

import (
	"net/http"

	"errors"

	"strconv"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func UserHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	username := mux.Vars(r)["ID"]

	ref, err := GetUserRef(store.DB, username)
	if err != nil {
		return rsp, err
	}

	user, err := GetUser(store, ref.Id)
	return handler.Response{
		Code: http.StatusOK,
		Data: user,
	}, nil
}

func LoggedInUserHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp, handler.StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Must be logged in to use this endpoint")}
	}

	usrRef := val.(model.Ref)
	user, err := GetUser(store, usrRef.Id)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: user,
	}, nil
}

func ImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	id := mux.Vars(r)["ID"]

	ref, err := GetImageRef(store.DB, id)
	if err != nil {
		return rsp, err
	}

	img, err := GetImage(store, ref.Id)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: img,
	}, nil
}

func TagHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	var err error
	var limit int
	id := mux.Vars(r)["ID"]

	params := r.URL.Query()
	l, ok := params["limit"]
	if ok {
		if len(l) == 1 {
			limit, err = strconv.Atoi(l[0])
			if err != nil {
				limit = 500
			}
		}
	}

	if limit == 0 {
		limit = 500
	}

	images, err := TaggedImages(store, id)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}
