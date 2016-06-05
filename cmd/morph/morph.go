package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/handler"
	"github.com/julienschmidt/httprouter"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port must be set")
	}

	router := httprouter.New()

	log.Printf("Serving at port:%s", port)
	router.GET("/", handler.IndexHandler)
	router.GET("/p/:img", handler.PictureHandler)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	router.ServeFiles("/content/*filepath", http.Dir("content/"))
	http.ListenAndServe(":"+port, router)
}
