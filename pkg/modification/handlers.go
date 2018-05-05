package modification

import (
	"errors"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/fokal/fokal-core/pkg/request"
	"github.com/fokal/fokal-core/pkg/retrieval"
	"github.com/fokal/fokal-core/pkg/stats"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
)

func FeatureImage(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)

	id := vars["ID"]

	imageRef, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = Feature(store.DB, imageRef.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}

func UnFeatureImage(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)

	id := vars["ID"]

	imageRef, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = UnFeature(store.DB, imageRef.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}

func PatchImage(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	vars := mux.Vars(r)

	id := vars["ID"]

	ref, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	req := new(request.PatchImageRequest)
	if err := binding.Bind(r, req); err != nil {
		return handler.Response{}, err
	}

	log.Printf("%+v\n", req)

	err = commitImagePatch(store.DB, ref, structs.Map(req))
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}

func PatchUser(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {

	user, ok := context.GetOk(r, "auth")
	if !ok {
		return handler.Response{}, handler.StatusError{Code: http.StatusUnauthorized, Err: errors.New("User is not logged in")}
	}

	ref := user.(model.Ref)

	req := new(request.PatchUserRequest)
	if err := binding.Bind(r, req); err != nil {
		return handler.Response{}, err
	}

	ref, err := retrieval.GetUserRef(store.DB, ref.Shortcode)
	if err != nil {
		return handler.Response{}, err
	}

	err = commitUserPatch(store.DB, ref, structs.Map(req))
	if err != nil {
		return handler.Response{}, err
	}
	return handler.Response{
		Code: http.StatusAccepted,
	}, nil

}

func DeleteImage(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	id := mux.Vars(r)["ID"]

	ref, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = deleteImage(store.DB, ref.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}

func DeleteUser(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	id := mux.Vars(r)["ID"]

	ref, err := retrieval.GetUserRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = deleteUser(store.DB, ref.Id)
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}

func DownloadHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	id := mux.Vars(r)["ID"]
	ref, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	err = stats.AddStat(store.DB, ref.Id, "download")
	if err != nil {
		return handler.Response{}, err
	}
	return handler.Response{Code: http.StatusAccepted}, nil

}
