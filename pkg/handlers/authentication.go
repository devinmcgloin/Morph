package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/sprioc/sprioc-core/pkg/authentication"
)

func GetToken(w http.ResponseWriter, r *http.Request) Response {
	decoder := json.NewDecoder(r.Body)

	creds := Credentials{}

	err := decoder.Decode(&creds)
	if err != nil {
		return Resp("Bad Request", http.StatusBadRequest)
	}

	log.Println(creds)

	valid, user, err := authentication.ValidateCredentialsByUserName(creds.Username, creds.Password)
	if err != nil {
		return Resp("Invalid Credentials", http.StatusUnauthorized)
	}
	if valid {
		token, err := authentication.CreateJWT(user)
		if err != nil {
			log.Println(err)
			return Resp("Invalid Credentials", http.StatusUnauthorized)
		}
		return Response{Code: 201, Data: structs.Map(tok{Token: token})}
	}

	return Resp("Invalid Credentials", http.StatusUnauthorized)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tok struct {
	Token string `json:"token"`
}
