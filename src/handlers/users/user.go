package users

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/devinmcgloin/sprioc/src/spriocError"
)

var mongo = store.ConnectStore()

func SignupHandler(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", http.StatusNotImplemented)
}

func AvatarUploadHander(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", http.StatusNotImplemented)

}

func formatSources(ID bson.ObjectId) model.ImgSource {
	const prefix = "https://images.sprioc.xyz/avatars/"
	var resourceBaseURL = prefix + ID.Hex()
	return model.ImgSource{
		Raw:    model.URL(resourceBaseURL),
		Large:  model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy"),
		Medium: model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max"),
		Small:  model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max"),
		Thumb:  model.URL(resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max"),
	}
}
