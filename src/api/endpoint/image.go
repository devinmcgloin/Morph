package endpoint

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/devinmcgloin/morph/src/api/SQL"
	"github.com/devinmcgloin/morph/src/views/common"
	"github.com/julienschmidt/httprouter"
)

func ImageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	IID, err := strconv.ParseUint(ps.ByName("IID"), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	img, err := SQL.GetImg(IID)
	if err != nil {
		log.Println(err)
		common.NotFound(w, r)
		return
	}

	err = decoder.Decode(&img, r.PostForm)

	log.Println(r.PostForm)

	log.Println(img)

	if err != nil {
		log.Println(err)
		common.SomethingsWrong(w, r, err)
		return
	}

	SQL.UpdateImg(img)

	newUrl := fmt.Sprintf("/i/%d/edit", IID)

	http.Redirect(w, r, newUrl, 302)

}
