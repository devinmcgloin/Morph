package handler

import (
	"errors"
	img "image"
	"net/http"

	"github.com/fokal/fokal-core/pkg/log"
	"github.com/fokal/fokal-core/pkg/services/image"
	"github.com/fokal/fokal-core/pkg/services/permission"

	"github.com/fokal/fokal-core/pkg/services/user"
	uuid "github.com/satori/go.uuid"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/fokal/fokal-core/pkg/request"
	"github.com/mholt/binding"
)

func RegisterHandlers(state *State, api *mux.Router, chain alice.Chain) {
	post := api.Methods("POST").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()
	get := api.Methods("GET").Subrouter()
	put := api.Methods("PUT").Subrouter()
	del := api.Methods("DELETE").Subrouter()
	patch := api.Methods("PATCH").Subrouter()

	post.Handle("/users", chain.Then(Handler{
		State: state,
		H:     CreateUser,
	}))
	opts.Handle("/users", chain.Then(Options("POST")))

	get.Handle("/users/{ID}", chain.Then(Handler{State: state, H: User}))
	opts.Handle("/users/{ID}", chain.Then(Options("GET")))

	put.Handle("/users/{ID}/featured",
		chain.Append(
			Middleware{State: state, M: Authenticate}.Handler,
			Permission{State: state,
				T:          permission.CanEdit,
				TargetType: permission.UserClass,
				M:          PermissionMiddle}.Handler).
			Then(Handler{State: state, H: FeatureUser}))
	del.Handle("/users/{ID}/featured",
		chain.Append(
			Middleware{State: state, M: Authenticate}.Handler,
			Permission{State: state,
				T:          permission.CanEdit,
				TargetType: permission.UserClass,
				M:          PermissionMiddle}.Handler).
			Then(Handler{State: state, H: UnFeatureUser}))

	opts.Handle("/users/{ID}/featured", chain.Then(Options("DELETE", "PUT")))

	put.Handle("/users/me/avatar", chain.Append(
		Middleware{
			State: state,
			M:     Authenticate,
		}.Handler).Then(Handler{
		State: state,
		H:     UploadAvatar,
	}))
	opts.Handle("/users/me/avatar", chain.Then(Options("PUT")))

	patch.Handle("/users/me",
		chain.Append(
			Middleware{State: state, M: Authenticate}.Handler,
			Permission{State: state,
				T:          permission.CanEdit,
				TargetType: permission.UserClass,
				M:          PermissionMiddle}.Handler).
			Then(Handler{State: state, H: PatchUser}))

	del.Handle("/users/me",
		chain.Append(
			Middleware{State: state, M: Authenticate}.Handler,
			Permission{State: state,
				T:          permission.CanEdit,
				TargetType: permission.UserClass,
				M:          PermissionMiddle}.Handler).
			Then(Handler{State: state, H: DeleteUser}))

	opts.Handle("/users/me", chain.Then(Options("PATCH", "DELETE")))
}

func CreateUser(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	ctx := r.Context()
	createRequest := &request.CreateUser{}

	if errs := binding.Bind(r, createRequest); errs != nil {
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
	name := claims["sub"].(string)

	user := &user.User{
		Email:    email,
		Name:     &name,
		Username: createRequest.Username,
	}
	err = s.UserService.CreateUser(ctx, user)
	if err != nil {
		return nil, &StatusError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}
	return &Response{
		Code: http.StatusOK,
	}, nil
}
func User(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	username := mux.Vars(r)["ID"]
	ctx := r.Context()
	user, err := s.UserService.UserByUsername(ctx, username)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}

	return &Response{Code: http.StatusOK, Data: user}, nil
}

func PatchUser(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	ctx := r.Context()
	userID := ctx.Value(log.UserIDKey).(uint64)
	patchRequest := &request.PatchUser{}

	if errs := binding.Bind(r, patchRequest); errs != nil {
		return nil, StatusError{
			Code: http.StatusBadRequest,
			Err:  errors.New("invalid Request: body missing required fields"),
		}
	}

	err := s.UserService.PatchUser(ctx, userID, *patchRequest)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}
	return &Response{Code: http.StatusAccepted}, nil
}

func DeleteUser(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	ctx := r.Context()
	userID := ctx.Value(log.UserIDKey).(uint64)

	err := s.UserService.DeleteUser(ctx, userID)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}
	return &Response{Code: http.StatusAccepted}, nil
}

func UploadAvatar(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	ctx := r.Context()
	log.WithContext(ctx).Debug("uploading avatar")
	userID := ctx.Value(log.UserIDKey).(uint64)
	id := uuid.NewV4()

	uploadedImage, _, err := img.Decode(r.Body)
	if err != nil {
		return nil, StatusError{
			Err:  errors.New("unable to read image body"),
			Code: http.StatusBadRequest}
	}

	if uploadedImage.Bounds().Dx() == 0 {
		return nil, StatusError{
			Err:  errors.New("cannot upload file with 0 bytes"),
			Code: http.StatusBadRequest}
	}

	err = s.StorageService.Avatar.UploadImage(ctx, uploadedImage, id.String())
	if err != nil {
		return nil, StatusError{
			Err:  errors.New("unable to upload image"),
			Code: http.StatusInternalServerError}
	}

	err = s.UserService.SetAvatarID(ctx, userID, id.String())
	if err != nil {
		return nil, StatusError{Code: http.StatusInternalServerError, Err: errors.New("Unable to update avatar id")}
	}

	return &Response{
		Code: http.StatusAccepted,
		Data: map[string]interface{}{"links": image.ImageSources(id.String(), "avatar")},
	}, nil
}

func FeatureUser(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	id := mux.Vars(r)["ID"]
	ctx := r.Context()

	log.WithContext(ctx).Debug("fetching user by username")
	user, err := s.UserService.UserByUsername(ctx, id)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}
	log.WithContext(ctx).Debug("setting user featured")
	err = s.UserService.Feature(ctx, user.ID)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}
	return &Response{Code: http.StatusAccepted}, nil
}

func UnFeatureUser(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	id := mux.Vars(r)["ID"]
	ctx := r.Context()

	user, err := s.UserService.UserByUsername(ctx, id)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}
	err = s.UserService.UnFeature(ctx, user.ID)
	if err != nil {
		return nil, StatusError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("unable to reach user service"),
		}
	}
	return &Response{Code: http.StatusAccepted}, nil
}
