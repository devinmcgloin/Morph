package publicView

import (
	"encoding/json"
	"net/http"

	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/devinmcgloin/sprioc/src/spriocError"
	"github.com/gorilla/mux"
)

func AlbumView(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userName := vars["username"]
	albumTitle := vars["Album"]

	album, err := mongo.GetAlbumCollectionView(userName, albumTitle)
	if err != nil {
		return spriocError.New(err, "Unable to get collection", 523)
	}

	usr, valid := session.GetUser(r)
	if valid {
		album.Auth = usr
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err = json.NewEncoder(w).Encode(album)

	if err != nil {
		return spriocError.New(err, "Unable to write JSON", 523)
	}
	return nil
}
