package handlers

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func GetCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func GetUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func GetImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
