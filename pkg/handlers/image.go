package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sprioc/sprioc-core/pkg/contentStorage"
	"github.com/sprioc/sprioc-core/pkg/metadata"
	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/store"
)

var decoder = schema.NewDecoder()

func GetImage(w http.ResponseWriter, r *http.Request) Response {

	img, _, err := getImage(r)
	if err != nil {
		return Resp("Image does not exist", http.StatusNotFound)
	}

	dat, err := json.Marshal(img)
	if err != nil {
		return Resp("Unable to write JSON", 523)
	}

	return Response{Code: http.StatusOK, Data: dat}
}

func UploadImage(w http.ResponseWriter, r *http.Request) Response {

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return Resp("Unauthorized Request, must be logged in to upload a new image.", http.StatusUnauthorized)
	}

	img := model.Image{
		ID:          bson.NewObjectId(),
		ShortCode:   mongo.GetNewImageShortCode(),
		User:        GetUserRef(user.ShortCode),
		PublishTime: time.Now(),
	}

	file, err := ioutil.ReadAll(r.Body)
	n := len(file)

	if n == 0 {
		return Resp("Cannot upload file with 0 bytes.", http.StatusBadRequest)
	}

	err = contentStorage.ProccessImage(file, n, img.ShortCode, "content")
	if err != nil {
		log.Println(err)
		return Resp(err.Error(), http.StatusBadRequest)
	}

	buf := bytes.NewBuffer(file)

	meta, err := metadata.GetMetadata(buf)
	if err != nil {
		return Resp("Error while reading image metadata", 500)
	}

	img.MetaData = meta

	img.Sources = formatImageSources(img.ShortCode)

	err = store.CreateImage(mongo, img)
	if err != nil {
		return Resp("Error while adding image to DB", 500)
	}

	err = store.ModifyUser(mongo, GetUserRef(user.ShortCode),
		bson.M{"$push": bson.M{"images": GetImageRef(img.ShortCode)}})
	if err != nil {
		log.Println(err)
		return Resp("Error while adding image to DB", 500)
	}

	return Response{Code: http.StatusOK}
}

func formatImageSources(shortcode string) model.ImgSource {
	const prefix = "https://images.sprioc.xyz/content/"
	var resourceBaseURL = prefix + shortcode
	return model.ImgSource{
		Raw:    resourceBaseURL,
		Large:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
		Medium: resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
		Small:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
		Thumb:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
	}
}

func FeatureImage(w http.ResponseWriter, r *http.Request) Response {
	return executeCommand(w, r, getImageInterface, store.FeatureImage)
}

func FavoriteImage(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getImageInterface, store.FavoriteImage)
}

func DeleteImage(w http.ResponseWriter, r *http.Request) Response {
	return executeCommand(w, r, getImageInterface, store.DeleteImage)
}

func UnFeatureImage(w http.ResponseWriter, r *http.Request) Response {
	return executeCommand(w, r, getImageInterface, store.UnFeatureImage)
}

func UnFavoriteImage(w http.ResponseWriter, r *http.Request) Response {
	return executeBiDirectCommand(w, r, getImageInterface, store.UnFavoriteImage)

}

func ModifyImage(w http.ResponseWriter, r *http.Request) Response {
	IID := mux.Vars(r)["IID"]
	ref := GetImageRef(IID)
	return executeCheckedModification(r, ref)
}
