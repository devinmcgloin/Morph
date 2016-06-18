package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/api/auth"
	"github.com/devinmcgloin/morph/src/api/endpoint"
	"github.com/devinmcgloin/morph/src/views/editView"
	"github.com/devinmcgloin/morph/src/views/publicView"
	"github.com/devinmcgloin/morph/src/views/settings"
	"github.com/gorilla/handlers"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func init() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port must be set")
	}

	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flag)

	gothic.Store = sessions.NewCookieStore(securecookie.GenerateRandomKey(16))

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "https://morph.devinmcgloin.com/auth/github/callback"),
	)

}

func main() {

	router := httprouter.New()
	port := os.Getenv("PORT")

	log.Printf("Serving at http://localhost:%s", port)

	// API POST ROUTES
	//TODO need to figure out api formatting and tokens.
	//TODO maybe these are the best routes for posting changes.

	router.POST("/api/upload", endpoint.UploadHandler)
	router.POST("/api/u/:UserName/edit", endpoint.UserHandler)
	router.POST("/api/i/:IID/edit", endpoint.ImageHandler)
	router.POST("/api/album/:AID/edit", endpoint.AlbumHandler)

	// CONTENT VIEW ROUTES
	router.GET("/", publicView.MostRecentView)
	router.GET("/i/:slug", publicView.FeatureImgView)
	router.GET("/tag/:tag", publicView.CollectionTagView)
	router.GET("/tag/:tag/:slug", publicView.CollectionTagFeatureView)
	router.GET("/album/:UserName/:Album", publicView.AlbumView)
	router.GET("/u/:UserName", publicView.UserProfileView)
	router.GET("/loc/:slug", publicView.LocationView)
	router.GET("/search/*query", publicView.SearchView)

	// CONTENT EDIT ROUTES
	router.GET("/i/:ImageShortTitle/edit", editView.FeatureImgEditView)
	router.GET("/album/:UserName/:slug/edit", editView.AlbumEditView)
	router.GET("/u/:UserName/edit", editView.UserProfileEditView)
	router.GET("/upload", editView.UploadView)

	// BACKEND MANAGE ROUTES
	router.GET("/login", publicView.UserLoginView)
	router.GET("/auth/:provider", auth.BeginAuthHandler)
	router.GET("/auth/:provider/callback", auth.UserLoginCallback)
	router.GET("/settings", settings.UserSettingsView)

	// ASSETS
	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)))

}
