package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/dbase"
	"github.com/devinmcgloin/morph/src/handler"
	"github.com/julienschmidt/httprouter"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port must be set")
	}

	router := httprouter.New()

	log.Printf("Serving at http://localhost:%s", port)
	router.GET("/", handler.IndexHandler)
	router.GET("/p/:p_id", handler.PictureHandler)
	router.GET("/admin/:type", handler.AdminHandler)
	router.POST("/upload", handler.UploadHandler)
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))
	dbase.SetDB()
	err := dbase.DB.Ping()
	if err != nil {
		log.Fatal("Database not connected. Ping Failed.")
	}
	http.ListenAndServe(":"+port, router)
}
