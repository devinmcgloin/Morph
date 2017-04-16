package handlers

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/core"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

func DeleteImage(w http.ResponseWriter, r *http.Request) rsp.Response {

	var user model.User
	val, _ := context.GetOk(r, "auth")

	user = val.(model.User)

	id := mux.Vars(r)["IID"]

	ref := refs.GetImageRef(id)

	return core.DeleteImage(user.GetRef(), ref)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.User
	val, _ := context.GetOk(r, "auth")

	user = val.(model.User)

	id := mux.Vars(r)["username"]

	ref := refs.GetUserRef(id)

	return core.DeleteUser(user.GetRef(), ref)
}
