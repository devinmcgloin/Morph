package social

import (
	"net/http"

	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/fokal/fokal-core/pkg/retrieval"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func FavoriteHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)

	id := vars["ID"]

	var usrRef model.Ref
	usr := context.Get(r, "auth")
	usrRef = usr.(model.Ref)

	imageRef, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = Favorite(store.DB, usrRef.Id, imageRef.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{Code: http.StatusAccepted}, nil
}

func UnFavoriteHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)

	id := vars["ID"]

	var usrRef model.Ref
	usr := context.Get(r, "auth")
	usrRef = usr.(model.Ref)

	imageRef, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = UnFavorite(store.DB, usrRef.Id, imageRef.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{Code: http.StatusAccepted}, nil
}

func FollowHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)

	id := vars["ID"]

	var usrRef model.Ref
	usr := context.Get(r, "auth")
	usrRef = usr.(model.Ref)

	followedRef, err := retrieval.GetUserRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = Follow(store.DB, usrRef.Id, followedRef.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{Code: http.StatusAccepted}, nil
}

func UnFollowHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)
	id := vars["ID"]

	var usrRef model.Ref
	usr := context.Get(r, "auth")
	usrRef = usr.(model.Ref)

	followedRef, err := retrieval.GetUserRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = UnFollow(store.DB, usrRef.Id, followedRef.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{Code: http.StatusAccepted}, nil
}
