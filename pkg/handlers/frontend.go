package handlers

import "net/http"

func LoadHTMLIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This will serve the react entry point."))
	// http.ServeFile(w, r, "dist/index.html")
}
