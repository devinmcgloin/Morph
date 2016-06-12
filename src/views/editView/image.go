package editView

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func FeatureImgEditView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	IID, err := strconv.Atoi(ps.ByName("IID"))
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}
	log.Printf("Accessing img:%d", uint64(IID))
	img, err := SQL.GetFeatureSingleImgView(uint64(IID))
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
