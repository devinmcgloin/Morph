package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/api"
	"github.com/devinmcgloin/morph/src/auth"
	"github.com/devinmcgloin/morph/src/dbase"
	"github.com/devinmcgloin/morph/src/handler"
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

	router.NotFound = http.HandlerFunc(handler.NotFound)

	log.Printf("Serving at http://localhost:%s", port)

	router.GET("/", handler.IndexHandler)
	router.GET("/p/:i_id", handler.PictureHandler)
	router.GET("/album/:album", handler.CategoryHandler)
	router.POST("/api/upload", api.UploadHandler)
	router.GET("/morph", handler.LoginDisplay)
	router.GET("/morph/:page", handler.AdminHandler)

	router.POST("/api/auth", auth.LoginHandler)

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	dbase.SetDB()

	http.ListenAndServe(":"+port, router)
}
