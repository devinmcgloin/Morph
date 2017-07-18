package security

import (
	"errors"
	"net/http"

	"log"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/request"
	"github.com/mholt/binding"
)

func LoginHandler(state *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	req := new(request.LoginRequest)
	if err := binding.Bind(r, req); err != nil {
		return handler.Response{}, err
	}

	log.Printf("%+v\n", req)

	valid, err := ValidateCredentials(state.DB, *req)
	if err != nil {
		return handler.Response{}, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
	}

	if valid {
		token, err := createJWT(state, model.Ref{Collection: model.Users, Shortcode: req.Username})
		if err != nil {
			return handler.Response{}, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
		}
		return handler.Response{Code: http.StatusOK, Data: map[string]string{"token": token}}, nil
	}

	return handler.Response{}, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
}

func PublicKeyHandler(state *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	key := make(map[string]string)
	key["554b5db484856bfa16e7da70a427dc4d9989678a"] = PublicKey
	return handler.Response{Code: http.StatusOK, Data: key}, nil
}
