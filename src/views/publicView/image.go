package publicView

import (
	"log"
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func FeatureImgView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ShortTitle := ps.ByName("shortcode")

	img, err := mongo.GetFeatureSingleImgView(ShortTitle)
	if err != nil {
		common.NotFound(w, r)
		return
	}

	t, err := common.StandardTemplate("templates/pages/imageView.tmpl")
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}

	log.Println(img)

	err = t.Execute(w, img)
	if err != nil {
		common.SomethingsWrong(w, r, err)
		return
	}
}
