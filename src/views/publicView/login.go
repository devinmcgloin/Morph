package publicView

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func UserLoginView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//TODO phase this out into a template
	fmt.Fprintf(w, "<p><a href='/auth/twitter?provider=twitter'>Click to log in with twitter</a></p>")
	fmt.Fprintf(w, "<p><a href='/auth/facebook?provider=facebook'>Click to log in with facebook</a></p>")
	fmt.Fprintf(w, "<p><a href='/auth/gplus?provider=gplus'>Click to log in with gplus</a></p>")
}
