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
