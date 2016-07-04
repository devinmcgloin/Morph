package handlers

import "net/http"

func LoadFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}
