package create

import (
	"image"
	"net/http"

	"errors"

	"io/ioutil"

	"bytes"

	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal-core/pkg/geo"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/metadata"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/fokal/fokal-core/pkg/retrieval"
	"github.com/fokal/fokal-core/pkg/tokens"
	"github.com/fokal/fokal-core/pkg/upload"
	"github.com/fokal/fokal-core/pkg/vision"
	"github.com/gorilla/context"
	uuid "github.com/satori/go.uuid"
)

func UserHandler(store *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	token, err := tokens.Parse(store, r)
	log.Println(token)
	if err == nil && token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return handler.Response{}, handler.StatusError{Code: http.StatusBadRequest, Err: errors.New(http.StatusText(http.StatusBadRequest))}
		}

		email := claims["email"].(string)
		name := claims["name"].(string)
		username := strings.Split(email, "@")[0]
		if domain, ok := claims["hd"]; ok {
			username = username + "." + domain.(string)
		}

		log.Printf("Creating new user: {Username: %s, Email: %s, Name: %s}", username, email, name)
		err = CommitUser(store.DB, username, email, name)
		if err != nil {
			return handler.Response{}, handler.StatusError{
				Code: http.StatusInternalServerError,
				Err:  errors.New("error while adding user to db")}
		}

		ref, err := retrieval.GetUserRef(store.DB, username)
		if err != nil {
			log.Println(err)
			return handler.Response{}, handler.StatusError{Code: http.StatusInternalServerError}
		}
		token, _ := tokens.Create(store, ref, email)
		return handler.Response{Code: http.StatusAccepted, Data: map[string]string{"token": token}}, nil
	} else {
		log.Println(token, err)
		return handler.Response{}, handler.StatusError{Code: http.StatusBadRequest, Err: errors.New("Token is invalid.")}
	}
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

	if uploadedImage.Bounds().Dx() <= 1500 || uploadedImage.Bounds().Dy() <= 1500 {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Cannot upload image smaller than 1500x1500"),
			Code: http.StatusBadRequest}
	}

	errChan := make(chan error, 3)
	metadataChan := make(chan model.ImageMetadata, 1)
	annotationsChan := make(chan vision.ImageResponse, 1)

	go metadata.GetMetadata(errChan, metadataChan, bytes.NewBuffer(file))

	go vision.AnnotateImage(errChan, annotationsChan, store.DB, store.Vision, uploadedImage)

	for i := 0; i < 2; i++ {
		err = <-errChan
		if err != nil {
			return handler.Response{}, err
		}

	}

	img.Metadata = <-metadataChan
	annotations := <-annotationsChan
	rotatedImage := metadata.NormalizeOrientatation(uploadedImage, img.Metadata.Orientation)
	img.Metadata.PixelXDimension = int64(rotatedImage.Bounds().Dx())
	img.Metadata.PixelYDimension = int64(rotatedImage.Bounds().Dy())

	go upload.ProccessImage(errChan, rotatedImage, format, img.Shortcode, "content")
	err = <-errChan
	if err != nil {
		return handler.Response{}, err
	}

	if !annotations.Safe {
		return handler.Response{}, handler.StatusError{
			Err:  errors.New("Image contains violent, medical or adult imagery."),
			Code: http.StatusBadRequest}
	}

	if img.Metadata.Location != nil {
		addr, err := geo.ReverseGeocode(store.Maps, img.Metadata.Location.Point)
		if err != nil {
			log.Println(err)
		}
		img.Metadata.Location.Description = &addr
	}

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
		Data: map[string]string{"link": ref.ToURL(store.Port, store.Local), "id": ref.Shortcode},
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

	uid := uuid.NewV4()

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

	go upload.ProccessImage(errChan, uploadedImage, format, uid.String(), "avatar")

	err = <-errChan
	if err != nil {
		return handler.Response{}, err
	}

	_, err = store.DB.Exec("UPDATE content.users set avatar_id = $1 where id = $2", uid.String(), user.Id)
	if err != nil {
		log.Println(err)
		return handler.Response{}, handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("Unable to update avatar id")}
	}

	return handler.Response{
		Code: http.StatusAccepted,
		Data: map[string]interface{}{"links": retrieval.ImageSources(uid.String(), "avatar")},
	}, nil
}
