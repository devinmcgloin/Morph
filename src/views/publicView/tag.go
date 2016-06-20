package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/morphError"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/gorilla/mux"
)

func CollectionTagView(w http.ResponseWriter, r *http.Request) error {

	tag := mux.Vars(r)["tag"]

	taggedImages, err := mongo.GetCollectionViewByTags(tag)
	if err != nil {
		return morphError.New(err, "Unable to get collection", 523)

	}

	if len(taggedImages.Images) == 0 {
		return morphError.New(err, "Collection was Empty", 404)

	}

	usr, valid := session.GetUser(r)
	if valid {
		taggedImages.Auth = usr
	}

	return common.ExecuteTemplate(w, r, "templates/public/tagView.tmpl", taggedImages)

}

func CollectionTagFeatureView(w http.ResponseWriter, r *http.Request) error {
	return morphError.New(nil, "Not Implemented", 404)
}
