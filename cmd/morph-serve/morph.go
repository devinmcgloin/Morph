package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/api"
	"github.com/devinmcgloin/morph/src/api/endpoint"

	"github.com/devinmcgloin/morph/src/viewHandler/publicView"
	"github.com/devinmcgloin/morph/src/viewHandler/secureView"
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

	// API POST ROUTES

	router.POST("/api/v0/upload", endpoint.UploadHandler)
	router.POST("/api/v0/users/:user", endpoint.UserHandler)
	router.POST("/api/v0/photos/:p_id", endpoint.ImageHandler)

	// CONTENT VIEW ROUTES
	router.GET("/", publicView.MostRecentView)
	router.GET("/i/:IID", publicView.FeatureImgView)
	router.GET("/tag/:tag", publicView.CollectionTagView)
	router.GET("/tag/:tag/:IID", publicView.CollectionTagFeatureView)
	router.GET("/album/:AID", publicView.AlbumCollectionView)
	router.GET("/u/:UID", publicView.UserProfileView)
	router.GET("/login", publicView.UserLoginView)
	router.GET("/loc/:LID", publicView.LocationView)
	router.GET("/search/*query", publicView.SearchView)

	// CONTENT EDIT ROUTES
	router.GET("/i/:IID/edit", publicView.FeatureImgEditView)
	router.GET("/album/:AID/edit", publicView.AlbumCollectionEditView)
	router.GET("/u/:UID/edit", publicView.UserProfileEditView)

	// BACKEND MANAGE ROUTES
	router.GET("/settings", secureView.UserSettingsView)

	// ASSETS
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	api.SetDB()

	log.Fatal(http.ListenAndServe(":"+port, router))

}
