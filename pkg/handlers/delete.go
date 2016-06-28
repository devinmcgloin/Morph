package handlers

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func DeleteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
func DeleteCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
func DeleteUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
