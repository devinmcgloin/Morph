package account

import (
	"net/http"

	"github.com/devinmcgloin/sprioc/src/spriocError"
)

// AlbumEditorView handles specific edits for a given album
func AlbumEditorView(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", 404)
}

// AlbumListView handles showing all the users albums
func AlbumListView(w http.ResponseWriter, r *http.Request) error {
	return spriocError.New(nil, "Not Implemented", 404)
}
