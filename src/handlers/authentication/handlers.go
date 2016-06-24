package authentication

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/auth"
	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/spriocError"
)

var mongo = store.ConnectStore()

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	creds := Credentials{}

	err := decoder.Decode(&creds)
	if err != nil {
		return spriocError.New(err, "Bad Request", http.StatusBadRequest)
	}

	log.Println(creds)

	valid, user, err := auth.ValidateCredentialsByUserName(creds.Username, creds.Password)
	if err != nil {
		return spriocError.New(err, "Invalid Credentials", http.StatusUnauthorized)
	}
	if valid {
		token, err := auth.CreateJWT(user)
		if err != nil {
			return spriocError.New(err, "Invalid Credentials", http.StatusUnauthorized)
		}
		w.Write([]byte(token))
		return nil
	}
	return spriocError.New(err, "Invalid Credentials", http.StatusUnauthorized)

}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func PublicKeyHandler(w http.ResponseWriter, r *http.Request) error {

	w.Write(auth.GetPublicKey())
	return nil
}
