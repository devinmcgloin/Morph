package retrieval

import (
	"net/http"

	"errors"

	"strconv"

	"log"

	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
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
	if err != nil {
		return rsp, err
	}

	var isFollowed bool
	val, ok := context.GetOk(r, "auth")
	if ok {
		LoggedInID := val.(model.Ref).Id
		err = store.DB.Get(&isFollowed, `
		SELECT TRUE
		FROM CONTENT.user_follows
		WHERE user_id = $1 AND followed_id = $2;
		`, LoggedInID, user.Id)
		if err == nil {
			user.FollowedByUser = &isFollowed
		}
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: user,
	}, nil
}

func UserImagesHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	username := mux.Vars(r)["ID"]

	ref, err := GetUserRef(store.DB, username)
	if err != nil {
		return rsp, err
	}

	user, err := GetUserImages(store, ref.Id)
	return handler.Response{
		Code: http.StatusOK,
		Data: user,
	}, nil
}

func UserFavoritesHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	username := mux.Vars(r)["ID"]

	ref, err := GetUserRef(store.DB, username)
	if err != nil {
		return rsp, err
	}

	user, err := GetUserFavorites(store, ref.Id)
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

func LoggedInUserImagesHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response

	val, ok := context.GetOk(r, "auth")
	if !ok {
		return rsp, handler.StatusError{
			Code: http.StatusUnauthorized,
			Err:  errors.New("Must be logged in to use this endpoint")}
	}

	usrRef := val.(model.Ref)
	images, err := GetUserImages(store, usrRef.Id)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
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

	var isFavorited bool
	val, ok := context.GetOk(r, "auth")
	if ok {
		LoggedInID := val.(model.Ref).Id
		err = store.DB.Get(&isFavorited, `
		SELECT TRUE
		FROM CONTENT.user_favorites
		WHERE user_id = $1 AND image_id = $2;
		`, LoggedInID, img.Id)
		if err == nil {
			img.FavoritedByUser = &isFavorited
		}
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

func RecentImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	var err error

	params := r.URL.Query()

	var limit int
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

	log.Printf("%d\n", limit)
	images, err := RecentImages(store, limit)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func FeaturedImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	var err error

	params := r.URL.Query()
	var limit int
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
	images, err := FeaturedImages(store, limit)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}

func TrendingImagesHander(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var rsp handler.Response
	var err error

	params := r.URL.Query()
	var limit int
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
	images, err := Trending(store, limit)
	if err != nil {
		return rsp, err
	}

	return handler.Response{
		Code: http.StatusOK,
		Data: images,
	}, nil
}
