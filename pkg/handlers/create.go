package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func CreateCollection(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func CreateImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

// TODO need to send more of this functionality to core
func CreateUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	decoder := json.NewDecoder(r.Body)

	var newUser map[string]string

	err := decoder.Decode(&newUser)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.CreateUser(newUser)
}
