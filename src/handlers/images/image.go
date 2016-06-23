package images

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/sprioc/src/api/contentStorage"
	"github.com/devinmcgloin/sprioc/src/api/metadata"
	"github.com/devinmcgloin/sprioc/src/handlers/references"

	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var mongo = store.ConnectStore()

var decoder = schema.NewDecoder()

func GetImage(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["ID"]

	img, err := mongo.GetImage(ref.GetImageRef(id))
	if err != nil {
		return spriocError.New(err, "Unable to get image", 523)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(img)
	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}

	return spriocError.New(nil, "Success", 200)
}

func UploadImage(w http.ResponseWriter, r *http.Request) error {

	var user model.User
	val, ok := context.GetOk(r, "auth")
	if ok {
		user = val.(model.User)
	} else {
		return spriocError.New(nil, "Unauthorized Request", http.StatusUnauthorized)
	}

	img := model.Image{
		ID:   mongo.GetNewImageID(),
		User: ref.GetUserRef(user.ID),
	}

	var file []byte
	n, err := r.Body.Read(file)
	if err != nil {
		return spriocError.New(err, "Unsupported Media Type", 415)
	}

	err = contentStorage.ProccessImage(file, n, img.ID)
	if err != nil {
		return spriocError.New(err, "Error while uploading image", 500)
	}

	buf := bytes.NewBuffer(file)

	meta, err := metadata.GetMetadata(buf)
	if err != nil {
		return spriocError.New(err, "Error while reading image metadata", 500)
	}

	img.MetaData = meta

	img.Sources = formatSources(img.ID)

	err = mongo.CreateImage(img)
	if err != nil {
		return spriocError.New(err, "Error while adding image to DB", 500)
	}

	return nil
}

func formatSources(ID bson.ObjectId) model.ImgSource {
	const prefix = "https://images.sprioc.xyz/content/"
	var resourceBaseURL = prefix + ID.Hex()
	return model.ImgSource{
		Raw:    model.URL(resourceBaseURL),
		Large:  model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy"),
		Medium: model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max"),
		Small:  model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max"),
		Thumb:  model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max"),
	}
}
