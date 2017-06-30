package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/core"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/rsp"
	"github.com/gorilla/context"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.Ref
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.Ref)
	} else {
		return rsp.Response{Message: "Unauthorized Request, must be logged in to upload a new image.",
			Code: http.StatusUnauthorized}
	}

	file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return rsp.Response{Message: "Unable to read image body.",
			Code: http.StatusBadRequest}
	}

	return core.UploadAvatar(user, file)
}

func UploadImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.Ref
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.Ref)
	} else {
		return rsp.Response{Message: "Unauthorized Request, must be logged in to upload a new image.",
			Code: http.StatusUnauthorized}
	}

	file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return rsp.Response{Message: "Unable to read image body.",
			Code: http.StatusBadRequest}
	}

	return core.UploadImage(user, file)
}
