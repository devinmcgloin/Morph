package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/context"
	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
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
	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return rsp.Response{Message: "Unauthorized Request, must be logged in to upload a new image.",
			Code: http.StatusUnauthorized}
	}

	file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return rsp.Response{Message: "Unable to read image body.",
			Code: http.StatusBadRequest}
	}

	go core.UploadImage(user, file)

	return rsp.Response{Code: http.StatusAccepted}
}
