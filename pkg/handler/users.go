package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fokal/fokal-core/pkg/services/user"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/fokal/fokal-core/pkg/request"
	"github.com/mholt/binding"
)

func RegisterHandlers(state *State, api *mux.Router, chain alice.Chain) {
	post := api.Methods("POST").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	post.Handle("/users", chain.Then(Handler{
		State: state,
		H:     CreateUser,
	}))
	opts.Handle("/users", chain.Then(Options("POST")))

}

func CreateUser(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	ctx := r.Context()
	createRequest := &request.CreateUser{}

	if errs := binding.Bind(r, createRequest); errs != nil {
		log.Error(errs.Error())
		return nil, StatusError{
			Code: http.StatusBadRequest,
			Err:  errors.New("invalid Request: body missing required fields"),
		}
	}

	token, err := s.AuthService.ParseToken(ctx, createRequest.Token)
	if err != nil {
		return nil, &StatusError{
			Code: http.StatusBadRequest,
			Err:  errors.New("token is invalid"),
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, StatusError{Code: http.StatusBadRequest, Err: errors.New("token is invalid")}
	}

	email := claims["email"].(string)
	name := claims["name"].(string)

	exists, err := s.UserService.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}

	if exists {
		return nil, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("User with email %s already exists", email),
		}
	}

	exists, err = s.UserService.ExistsByUsername(ctx, createRequest.Username)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}

	if exists {
		return nil, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("user with username %s already exists", createRequest.Username),
		}
	}

	user := &user.User{
		Email:    email,
		Name:     &name,
		Username: createRequest.Username,
	}
	err = s.UserService.CreateUser(ctx, user)
	if err != nil {
		return nil, &StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to create user"),
		}
	}
	return &Response{
		Code: http.StatusOK,
	}, nil
}

// func PatchUser(s *State, w http.ResponseWriter, r *http.Request) (Response, error)     {}
// func DeleteUser(s *State, w http.ResponseWriter, r *http.Request) (Response, error)    {}
// func UploadAvatar(s *State, w http.ResponseWriter, r *http.Request) (Response, error)  {}
// func FeatureUser(s *State, w http.ResponseWriter, r *http.Request) (Response, error)   {}
// func UnFeatureUser(s *State, w http.ResponseWriter, r *http.Request) (Response, error) {}
