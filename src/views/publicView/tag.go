package publicView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func CollectionTagView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	tag := ps.ByName("tag")

	log.Printf("Accessing tag:%s", tag)
	taggedImages, err := SQL.GetCollectionViewByTag(tag)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	if len(taggedImages.Images) == 0 {
		common.NotFound(w, r)
		return
	}

	t, err := common.StandardTemplate("templates/pages/tagView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, taggedImages)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

}

func CollectionTagFeatureView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {}
