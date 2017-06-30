package handlers

import (
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/core"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/rsp"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func DeleteImage(w http.ResponseWriter, r *http.Request) rsp.Response {

	var user model.Ref
	val, _ := context.GetOk(r, "auth")

	user = val.(model.Ref)

	id := mux.Vars(r)["IID"]
	img, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}
	return core.DeleteImage(user, img)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	var user model.Ref
	val, _ := context.GetOk(r, "auth")

	user = val.(model.Ref)

	id := mux.Vars(r)["username"]

	usr, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.DeleteUser(user, usr)
}
