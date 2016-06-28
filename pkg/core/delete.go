package core

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func DeleteImage(requestFrom model.User, image model.DBRef) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
func DeleteCollection(requestFrom model.User, collection model.DBRef) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
func DeleteUser(requestFrom model.User, user model.DBRef) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
