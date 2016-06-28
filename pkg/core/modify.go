package core

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"gopkg.in/mgo.v2/bson"
)

func Modify(ref model.DBRef, changes bson.M) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented}
}
