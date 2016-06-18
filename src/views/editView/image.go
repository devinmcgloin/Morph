package editView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func FeatureImgEditView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	imageShortTitle := ps.ByName("ImageShortTitle")

	log.Printf("Accessing img:%s", imageShortTitle)
	img, err := mongo.GetFeatureSingleImgView(imageShortTitle)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	t, err := common.StandardTemplate("templates/pages/imageEditView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	err = t.Execute(w, img)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}
}
