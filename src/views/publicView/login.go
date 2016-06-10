package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func UserLoginView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	common.RenderStatic(w, r, "templates/pages/uploadView.tmpl")
}
