package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/core"
	"github.com/devinmcgloin/fokal/pkg/rsp"
)

func CreateUser(w http.ResponseWriter, r *http.Request) rsp.Response {
	decoder := json.NewDecoder(r.Body)

	var newUser map[string]string

	err := decoder.Decode(&newUser)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	return core.CreateUser(newUser)
}
