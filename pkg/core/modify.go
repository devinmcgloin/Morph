package core

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"github.com/sprioc/sprioc-core/pkg/store"
	"gopkg.in/mgo.v2/bson"
)

func Modify(ref model.DBRef, changes bson.M) rsp.Response {
	err := store.Modify(ref, changes)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

// TODO would really like to lock this down more and do more content validation.
func ModifySecure(user model.User, target model.DBRef, changes bson.M) rsp.Response {

	resp := VerifyChanges(user, target, changes)
	if !resp.Ok() {
		return resp
	}

	err := store.Modify(target, changes)
	if err != nil {
		return rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
	}

	return rsp.Response{Code: 200}
}
