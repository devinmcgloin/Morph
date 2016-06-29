package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/sprioc/sprioc-core/pkg/core"
	"github.com/sprioc/sprioc-core/pkg/rsp"
)

func GetToken(w http.ResponseWriter, r *http.Request) rsp.Response {
	decoder := json.NewDecoder(r.Body)

	creds := Credentials{}

	err := decoder.Decode(&creds)
	if err != nil {
		return rsp.Response{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	valid, user, err := core.ValidateCredentialsByUserName(creds.Username, creds.Password)
	if err != nil {
		return rsp.Response{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
	}
	if valid {
		token, resp := core.CreateJWT(user)
		if !resp.Ok() {
			log.Println(resp)
			return resp
		}
		return rsp.Response{Code: 201, Data: structs.Map(tok{Token: token})}
	}

	return rsp.Response{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tok struct {
	Token string `json:"token"`
}
