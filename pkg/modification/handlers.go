package modification

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/request"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
)

//func FeatureImage(w http.ResponseWriter, r *http.Request) rsp.Response {
//	var usrRef model.Ref
//	vars := mux.Vars(r)
//
//	id := vars["IID"]
//
//	usr, ok := context.GetOk(r, "auth")
//	if ok {
//		usrRef = usr.(model.Ref)
//	} else {
//		return rsp.Response{
//			Message: "Unauthorized Request, must be logged in to modify an image",
//			Code:    http.StatusUnauthorized,
//		}
//	}
//
//	imageRef, resp := core.GetImageRef(id)
//	if !resp.Ok() {
//		return resp
//	}
//
//	return core.FeatureImage(usrRef, imageRef)
//}
//
//func UnFeatureImage(w http.ResponseWriter, r *http.Request) rsp.Response {
//	var usrRef model.Ref
//	vars := mux.Vars(r)
//
//	id := vars["IID"]
//
//	usr, ok := context.GetOk(r, "auth")
//	if ok {
//		usrRef = usr.(model.Ref)
//	} else {
//		return rsp.Response{
//			Message: "Unauthorized Request, must be logged in to modify an image",
//			Code:    http.StatusUnauthorized,
//		}
//	}
//
//	imageRef, resp := core.GetImageRef(id)
//	if !resp.Ok() {
//		return resp
//	}
//
//	return core.UnFeatureImage(usrRef, imageRef)
//}

func PatchImage(w http.ResponseWriter, r *http.Request) rsp.Response {

	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	image, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	decoder := json.NewDecoder(r.Body)

	var request map[string]interface{}

	err := decoder.Decode(&request)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.PatchImage(usrRef, image, request)
}

func PatchUser(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {

	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["username"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return nil, handler.StatusError{
			Err:  errors.New("Unauthorized Request, must be logged in to modify an image"),
			Code: http.StatusUnauthorized,
		}
	}

	req := new(request.PatchUserRequest)
	if err := binding.Bind(r, req); err != nil {
		return nil, err
	}

}

//
//func DeleteImage(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
//	var user model.Ref
//	val, _ := context.GetOk(r, "auth")
//
//	user = val.(model.Ref)
//
//	id := mux.Vars(r)["IID"]
//
//}
//
//func DeleteUser(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
//	var user model.Ref
//	val, _ := context.GetOk(r, "auth")
//
//	user = val.(model.Ref)
//
//	id := mux.Vars(r)["username"]
//
//}
