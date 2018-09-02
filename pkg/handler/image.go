package handler

import (
	"errors"
	img "image"
	"net/http"

	"github.com/fokal/fokal-core/pkg/log"
	"github.com/fokal/fokal-core/pkg/services/image"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func registerImageHandlers(state *State, api *mux.Router, chain alice.Chain) {
	post := api.Methods("POST").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()
	//get := api.Methods("GET").Subrouter()
	//put := api.Methods("PUT").Subrouter()
	//del := api.Methods("DELETE").Subrouter()
	//patch := api.Methods("PATCH").Subrouter()

	post.Handle("/images", chain.Append(
		Middleware{
			State: state,
			M:     Authenticate,
		}.Handler).Then(Handler{
		State: state,
		H:     CreateImage,
	}))
	opts.Handle("/images", chain.Then(Options("POST")))
}

func CreateImage(s *State, w http.ResponseWriter, r *http.Request) (*Response, error) {
	ctx := r.Context()
	log.WithContext(ctx).Debug("uploading new image")
	userID := ctx.Value(log.UserIDKey).(uint64)

	shortcode, err := s.ImageService.NextShortcode(ctx)
	if err != nil {
		return nil, StatusError{
			Err:  errors.New("unable to generate new shortcode"),
			Code: http.StatusBadRequest}
	}
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

	err = s.StorageService.Content.UploadImage(ctx, uploadedImage, shortcode)
	if err != nil {
		return nil, StatusError{
			Err:  errors.New("unable to upload image"),
			Code: http.StatusInternalServerError}
	}

	err = s.ImageService.CreateImage(ctx, &image.Image{
		Shortcode: shortcode,
		UserID:    userID,
	})
	if err != nil {
		return nil, StatusError{
			Err:  errors.New("unable to upload image"),
			Code: http.StatusInternalServerError}
	}

	return &Response{
		Code: http.StatusAccepted,
	}, nil
}
