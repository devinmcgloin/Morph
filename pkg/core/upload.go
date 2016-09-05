package core

import (
	"errors"
	"log"
	"net/http"

	"github.com/sprioc/composer/pkg/contentStorage"
	"github.com/sprioc/composer/pkg/metadata"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

var clarifaijobs chan string

func init() {
	clarifaijobs = make(chan string)
	metadata.SetupClarifai(clarifaijobs)
	metadata.Start()
}

// TODO NEED TO ABSTRACT THIS FURTHER

// func UploadImage(user model.User, file []byte) rsp.Response {
//
// 	img := model.Image{
// 		ID:          bson.NewObjectId(),
// 		ShortCode:   store.GetNewImageShortCode(),
// 		Owner:       refs.GetUserRef(user.ShortCode),
// 		PublishTime: time.Now(),
// 	}
//
// 	n := len(file)
//
// 	if n == 0 {
// 		return rsp.Response{Message: "Cannot upload file with 0 bytes.", Code: http.StatusBadRequest}
// 	}
//
// 	err := contentStorage.ProccessImage(file, n, img.ShortCode, "content")
// 	if err != nil {
// 		log.Println(err)
// 		return rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
// 	}
//
// 	buf := bytes.NewBuffer(file)
//
// 	img.MetaData = metadata.GetMetadata(buf)
//
// 	img.Sources = formatSources(img.ShortCode, "content")
//
// 	err = store.Create("images", img)
// 	if err != nil {
// 		return rsp.Response{Message: "Error while adding image to DB", Code: http.StatusInternalServerError}
// 	}
//
// 	imgRef := refs.GetImageRef(img.ShortCode)
// 	resp := Modify(refs.GetUserRef(user.ShortCode),
// 		bson.M{"$push": bson.M{"images": imgRef}})
// 	if !resp.Ok() {
// 		return rsp.Response{Message: "Error while adding image to DB", Code: http.StatusInternalServerError}
// 	}
//
// 	go metadata.SetLocation(img.MetaData.Location)
// 	go func() { clarifaijobs <- img.ShortCode }()
//
// 	return rsp.Response{Code: http.StatusAccepted, Data: map[string]string{"link": refs.GetURL(imgRef)}}
// }

func UploadAvatar(user model.User, file []byte) rsp.Response {
	n := len(file)

	if n == 0 {
		return rsp.Response{Message: "Cannot upload file with 0 bytes.", Code: http.StatusBadRequest}
	}

	// REVIEW check that you can overwrite on aws
	err := contentStorage.ProccessImage(file, n, user.ShortCode, "avatar")
	if err != nil {
		log.Println(err)
		return rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
	}

	sources := formatSources(user.ShortCode, "avatars")

	err = setAvatar(refs.GetUserRef(user.ShortCode), sources)
	if err != nil {
		return rsp.Response{Message: "Unable to add image", Code: http.StatusInternalServerError}
	}
	return rsp.Response{Code: http.StatusAccepted}
}

func formatSources(shortcode, location string) model.ImgSource {
	var prefix = "https://images.sprioc.xyz/" + location + "/"
	var resourceBaseURL = prefix + shortcode
	return model.ImgSource{
		Raw:    resourceBaseURL,
		Large:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
		Medium: resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
		Small:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
		Thumb:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
	}
}

// TODO Implement avatar upload
func setAvatar(user model.Ref, source model.ImgSource) error {
	return errors.New("NOT IMPLEMENTED")
}
