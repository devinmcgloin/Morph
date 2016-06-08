package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/api"
	"github.com/devinmcgloin/morph/src/api/endpoint"
	"github.com/julienschmidt/httprouter"
)

func main() {

	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flag)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port must be set")
	}

	router := httprouter.New()

	log.Printf("Serving at http://localhost:%s", port)

	router.POST("/api/v0/upload", endpoint.UploadHandler)
	router.POST("/api/v0/auth", endpoint.LoginHandler)
	router.POST("/api/v0/users/:user", endpoint.UserHandler)
	router.POST("/api/v0/photos/:p_id", endpoint.ImageHandler)

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	api.SetDB()

	log.Fatal(http.ListenAndServe(":"+port, router))

}
