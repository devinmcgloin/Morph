package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/core"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/rsp"
)

func GetToken(w http.ResponseWriter, r *http.Request) rsp.Response {
	decoder := json.NewDecoder(r.Body)

	var creds = make(map[string]string)
	var username, password string
	var ok bool

	err := decoder.Decode(&creds)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if username, ok = creds["username"]; !ok {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if password, ok = creds["password"]; !ok {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	valid, resp := core.ValidateCredentialsByUserName(username, password)
	if !resp.Ok() {
		return rsp.Response{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
	}

	if valid {
		token, resp := core.CreateJWT(model.Ref{Collection: model.Users, Shortcode: username})
		if !resp.Ok() {
			log.Println(resp)
			return resp
		}
		return rsp.Response{Code: 201, Data: tok{Token: token}}
	}

	return rsp.Response{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
}

type tok struct {
	Token string `json:"token"`
}
