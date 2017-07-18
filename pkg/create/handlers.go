package create

import (
	"image"
	"net/http"

	"errors"

	"io/ioutil"

	"bytes"

	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/metadata"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/request"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/devinmcgloin/fokal/pkg/upload"
	"github.com/devinmcgloin/fokal/pkg/vision"
	"github.com/gorilla/context"
	"github.com/mholt/binding"
)

func UserHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	req := new(request.CreateUserRequest)
	if err := binding.Bind(r, req); err != nil {
		return handler.Response{}, err
	}

	err := validateUser(store.DB, req)
	if err != nil {
		return handler.Response{}, err
	}

	securePassword, salt, err := security.GenerateSaltPass(req.Password)
	if err != nil {
		return handler.Response{}, err
	}

	usr := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: securePassword,
		Salt:     salt,
	}

	err = commitUser(store.DB, usr)
	if err != nil {
		return handler.Response{}, err
	}

	ref := model.Ref{Collection: model.Users, Shortcode: usr.Username}
	return handler.Response{
		Code: http.StatusAccepted,
		Data: map[string]string{"link": ref.ToURL(store.Port, store.Local)},
	}, nil
}

func ImageHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var user model.Ref
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.Ref)
	} else {
		return handler.Response{}, handler.StatusError{Code: http.StatusUnauthorized}
	}

	file, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Unable to read image body."),
			Code: http.StatusBadRequest}
	}

	uploadedImage, format, err := image.Decode(bytes.NewBuffer(file))
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Unable to read image body."),
			Code: http.StatusBadRequest}
	}

	sc, err := retrieval.GenerateSC(store.DB, model.Images)
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Unable to generate new shortcode"),
			Code: http.StatusInternalServerError}
	}

	img := model.Image{
		Shortcode: sc,
		UserId:    user.Id,
	}

	if uploadedImage.Bounds().Dx() == 0 {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Cannot upload file with 0 bytes."),
			Code: http.StatusBadRequest}
	}

	errChan := make(chan error, 3)
	metadataChan := make(chan model.ImageMetadata, 1)
	annotationsChan := make(chan vision.ImageResponse, 1)

	go upload.ProccessImage(errChan, uploadedImage, format, img.Shortcode, "content")

	go metadata.GetMetadata(errChan, metadataChan, bytes.NewBuffer(file))

	go vision.AnnotateImage(errChan, annotationsChan, store.DB, store.Vision, uploadedImage)

	for i := 0; i < 3; i++ {
		select {
		case err := <-errChan:
			if err != nil {
				return handler.Response{}, err
			}
		}
	}

	img.Metadata = <-metadataChan
	annotations := <-annotationsChan
	img.Labels = annotations.Labels
	img.Landmarks = annotations.Landmark
	img.Colors = annotations.ColorProperties

	err = commitImage(store.DB, img)
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Error while adding image to DB"),
			Code: http.StatusInternalServerError}
	}

	ref := model.Ref{Collection: model.Images, Shortcode: img.Shortcode}
	return handler.Response{
		Code: http.StatusAccepted,
		Data: map[string]string{"link": ref.ToURL(store.Port, store.Local)},
	}, nil

}

func AvatarHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	var user model.Ref
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.Ref)
	} else {
		return handler.Response{}, handler.StatusError{Code: http.StatusUnauthorized}
	}

	uploadedImage, format, err := image.Decode(r.Body)
	if err != nil {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Unable to read image body."),
			Code: http.StatusBadRequest}
	}

	if uploadedImage.Bounds().Dx() == 0 {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Cannot upload file with 0 bytes."),
			Code: http.StatusBadRequest}
	}

	errChan := make(chan error, 1)

	go upload.ProccessImage(errChan, uploadedImage, format, user.Shortcode, "avatar")
	select {
	case err := <-errChan:
		if err != nil {
			return handler.Response{}, err
		}
	}

	return handler.Response{
		Code: http.StatusAccepted,
	}, nil
}
