package handlers

import "net/http"

func LoadHTMLIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "dist/index.html")
}
