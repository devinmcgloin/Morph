package publicView

import (
	"net/http"

	"github.com/devinmcgloin/morph/src/views/common"
)

func UserLoginView(w http.ResponseWriter, r *http.Request) {

	common.ExecuteTemplate(w, r, "templates/pages/loginView.tmpl", nil)
}
