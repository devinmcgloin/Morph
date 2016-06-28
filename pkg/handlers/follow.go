package handlers

import (
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func UnFollowCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func FollowCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
func UnFollowUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func FollowUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
