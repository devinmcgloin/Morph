package handler

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/fokal/fokal-core/pkg/domain"
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

	exists, err := s.UserService.ExistsByEmail(ctx, createRequest.Email)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}

	if exists {
		return nil, StatusError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("User with email %s already exists", createRequest.Email),
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

	valid, _, err := s.AuthService.VerifyToken(ctx, createRequest.Token)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	if !valid {
		return nil, &StatusError{
			Code: http.StatusBadRequest,
			Err:  errors.New("token is invalid"),
		}
	}

	user := &domain.User{
		Email:    createRequest.Email,
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
