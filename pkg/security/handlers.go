package security

import (
	"errors"
	"net/http"

	"log"

	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/request"
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
		jwt, err := createJWT(state, model.Ref{Collection: model.Users, Shortcode: req.Username})
		if err != nil {
			return handler.Response{}, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
		}
		return handler.Response{Code: http.StatusOK, Data: map[string]string{"jwt": jwt}}, nil
	}

	return handler.Response{}, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
}

func PublicKeyHandler(state *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	key := make(map[string]string)
	key[keyHash] = PublicKey
	return handler.Response{Code: http.StatusOK, Data: key}, nil
}

func RefreshHandler(state *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	user, err := verifyJWT(state, r)
	if err != nil {
		return handler.Response{}, handler.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	jwt, err := createJWT(state, user)
	if err != nil {
		return handler.Response{}, handler.StatusError{Code: http.StatusBadRequest, Err: err}
	}
	return handler.Response{Code: http.StatusOK, Data: map[string]string{"token": jwt}}, nil
}
