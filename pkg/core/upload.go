package core

import (
	"bytes"
	"errors"
	"log"
	"net/http"

	"github.com/sprioc/composer/pkg/contentStorage"
	"github.com/sprioc/composer/pkg/metadata"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/rsp"
	"github.com/sprioc/composer/pkg/sql"
)

// TODO NEED TO ABSTRACT THIS FURTHER

func UploadImage(user model.Ref, file []byte) rsp.Response {

	var err error
	sc, err := sql.GenerateSC(model.Images)
	if err != nil {
		log.Println(err)
		return rsp.Response{Message: "Unable to generate new shortcode", Code: http.StatusInternalServerError}
	}

	img := model.Image{
		Shortcode: sc,
		OwnerId:   user.Id,
	}

	n := len(file)

	if n == 0 {
		return rsp.Response{Message: "Cannot upload file with 0 bytes.", Code: http.StatusBadRequest}
	}

	err = contentStorage.ProccessImage(file, n, img.Shortcode, "content")
	if err != nil {
		log.Println(err)
		return rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
	}

	buf := bytes.NewBuffer(file)

	metadata.GetMetadata(buf, &img.Metadata)

	annotations, err := metadata.AnnotateImage(bytes.NewBuffer(file))
	if err != nil {
		log.Println(err)
	} else {
		img.Labels = annotations.Labels
		img.Landmarks = annotations.Landmark
		img.Colors = annotations.ColorProperties
	}

	err = sql.CreateImage(img)
	if err != nil {
		return rsp.Response{Message: "Error while adding image to DB", Code: http.StatusInternalServerError}
	}

	return rsp.Response{Code: http.StatusAccepted, Data: map[string]string{"link": ""}}
}

func UploadAvatar(userRef model.Ref, file []byte) rsp.Response {
	n := len(file)

	if n == 0 {
		return rsp.Response{Message: "Cannot upload file with 0 bytes.", Code: http.StatusBadRequest}
	}

	user, err := sql.GetUser(userRef.Id)
	if err != nil {
		return rsp.Response{Code: http.StatusNotFound}
	}

	// REVIEW check that you can overwrite on aws
	err = contentStorage.ProccessImage(file, n, user.Username, "avatar")
	if err != nil {
		log.Println(err)
		return rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
	}

	return rsp.Response{Code: http.StatusAccepted}
}

// TODO Implement avatar upload
func setAvatar(user model.Ref, source model.ImageSource) error {
	return errors.New("NOT IMPLEMENTED")
}
