package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sprioc/sprioc-core/pkg/contentStorage"
	"github.com/sprioc/sprioc-core/pkg/metadata"
	"github.com/sprioc/sprioc-core/pkg/model"
)

var decoder = schema.NewDecoder()

func GetImage(w http.ResponseWriter, r *http.Request) Response {

	id := mux.Vars(r)["ID"]

	img, err := mongo.GetImage(GetImageRef(id))
	if err != nil {
		return Resp("Unable to get image", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
		return Resp("Unauthorized Request", http.StatusUnauthorized)
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

	err = mongo.CreateImage(img)
	if err != nil {
		return Resp("Error while adding image to DB", 500)
	}

	err = mongo.ModifyUser(GetUserRef(user.ShortCode),
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

func FeatureHandler(w http.ResponseWriter, r *http.Request) Response {

}

func FavoriteHandler(w http.ResponseWriter, r *http.Request) Response {

}

func DeleteImage(w http.ResponseWriter, r *http.Request) Response {
	id := mux.Vars(r)["ID"]
	img, err := mongo.GetImage(GetImageRef(id))
	if err != nil {
		return Resp("Unable to get image", http.StatusNotFound)
	}

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return Resp("Unauthorized Request", http.StatusUnauthorized)
	}

	if strings.Compare(user.ShortCode, img.User.Shortcode) != 0 {
		return Resp("Unauthorized Request", http.StatusUnauthorized)
	}

	err = mongo.DeleteImage(GetImageRef(id))
	if err != nil {
		return Resp("Internal Server Error", http.StatusInternalServerError)
	}
	return Response{Code: http.StatusAccepted}
}

func UnFeatureHandler(w http.ResponseWriter, r *http.Request) Response {

}

func UnFavoriteHandler(w http.ResponseWriter, r *http.Request) Response {

}

func ChangeHandler(w http.ResponseWriter, r *http.Request) Response {
	id := mux.Vars(r)["ID"]
	img, err := mongo.GetImage(GetImageRef(id))
	if err != nil {
		return Resp("Unable to get image", http.StatusNotFound)
	}

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return Resp("Unauthorized Request", http.StatusUnauthorized)
	}

	if strings.Compare(user.ShortCode, img.User.Shortcode) != 0 {
		return Resp("Unauthorized Request", http.StatusUnauthorized)
	}

	var m bson.M

	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Println(err)
		return Resp("Malformed Request", http.StatusBadRequest)
	}

	log.Println(m)

	err = mongo.ModifyImage(GetImageRef(id), m)
	if err != nil {
		return Resp(err.Error(), http.StatusBadRequest)
	}

	return Response{Code: 200}
}
