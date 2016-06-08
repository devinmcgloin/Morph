package secureView

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func UserDashboardView() (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {}
