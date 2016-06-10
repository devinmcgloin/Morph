package editView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func UploadView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	common.RenderStatic(w, r, "templates/pages/uploadView.tmpl")
}
