package modification

import (
	"net/http"

	"log"

	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/request"
	"github.com/fokal/fokal/pkg/retrieval"
	"github.com/fatih/structs"
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
	log.Println("Inside PatchImage")

	id := vars["ID"]

	ref, err := retrieval.GetImageRef(store.DB, id)
	if err != nil {
		return handler.Response{}, err
	}

	req := new(request.PatchImageRequest)
	if err := binding.Bind(r, req); err != nil {
		return handler.Response{}, err
	}

	err = commitImagePatch(store.DB, ref, structs.Map(req))
	if err != nil {
		return handler.Response{}, err
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}

func PatchUser(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {

	vars := mux.Vars(r)

	id := vars["ID"]

	req := new(request.PatchUserRequest)
	if err := binding.Bind(r, req); err != nil {
		return handler.Response{}, err
	}

	ref, err := retrieval.GetUserRef(store.DB, id)
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
