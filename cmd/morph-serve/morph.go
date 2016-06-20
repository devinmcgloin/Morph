package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/morph/src/api/auth"
	"github.com/devinmcgloin/morph/src/api/endpoint"
	"github.com/devinmcgloin/morph/src/api/session"
	"github.com/devinmcgloin/morph/src/views/account"
	"github.com/devinmcgloin/morph/src/views/publicView"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
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
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK_URL")),
	)

}

func main() {

	router := mux.NewRouter()
	port := os.Getenv("PORT")

	log.Printf("Serving at http://localhost:%s", port)

	// API POST ROUTES
	// TODO need to figure out api formatting and tokens.
	// TODO make shortcodes a primary field in mongo and index by those.

	router.HandleFunc("/api/upload", setAuthContext(endpoint.UploadHandler)).Methods("POST")
	router.HandleFunc("/api/u/{username}/edit", setAuthContext(endpoint.UserHandler)).Methods("POST")
	router.HandleFunc("/api/i/{shortcode}/edit", setAuthContext(endpoint.ImageHandler)).Methods("POST")
	router.HandleFunc("/api/album/{username}/{shortcode}/edit", setAuthContext(endpoint.AlbumHandler)).Methods("POST")

	// CONTENT VIEW ROUTES
	router.HandleFunc("/", setAuthContext(publicView.MostRecentView)).Methods("GET")
	router.HandleFunc("/i/{shortcode}/", setAuthContext(publicView.FeatureImgView)).Methods("GET")
	router.HandleFunc("/tag/{tag}/", setAuthContext(publicView.CollectionTagView)).Methods("GET")
	router.HandleFunc("/tag/{tag}/{shortcode}/", setAuthContext(publicView.CollectionTagFeatureView)).Methods("GET")
	router.HandleFunc("/album/{username}/{shortcode}/", setAuthContext(publicView.AlbumView)).Methods("GET")
	router.HandleFunc("/u/{username}/", setAuthContext(publicView.UserProfileView)).Methods("GET")
	router.HandleFunc("/loc/", setAuthContext(publicView.LocationView)).Methods("GET") //TODO not sure about shortcodes for locations
	router.HandleFunc("/search/", setAuthContext(publicView.SearchView)).Methods("GET")

	// CONTENT EDIT ROUTES
	router.HandleFunc("/account/images/", secureAuthContext(account.ImageEditorView)).Methods("GET")
	router.HandleFunc("/account/albums/", secureAuthContext(account.AlbumListView)).Methods("GET")
	router.HandleFunc("/account/albums/{shortcode}/", secureAuthContext(account.AlbumEditorView)).Methods("GET")
	router.HandleFunc("/upload/", secureAuthContext(account.UploadView)).Methods("GET")
	router.HandleFunc("/account/", secureAuthContext(account.UserSettingsView)).Methods("GET")

	// BACKEND MANAGE ROUTES
	router.HandleFunc("/login/", setAuthContext(account.UserLoginView)).Methods("GET")
	router.HandleFunc("/logout/", setAuthContext(account.UserLogoutView)).Methods("GET")
	router.HandleFunc("/auth/{provider}", setAuthContext(auth.BeginAuthHandler)).Methods("GET")
	router.HandleFunc("/auth/{provider}/callback", setAuthContext(auth.UserLoginCallback)).Methods("GET")

	// ASSETS
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)))

}

func setAuthContext(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedIn, user := auth.CheckUser(r)
		if loggedIn {
			log.Printf("User: %s", user.UserName)
			w = auth.RenewCookie(w, r)
			session.SetUser(r, user)
		}

		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	})
}

func secureAuthContext(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggedIn, user := auth.CheckUser(r)
		if loggedIn {
			log.Printf("User: %s", user.UserName)

			w = auth.RenewCookie(w, r)
			session.SetUser(r, user)
			err := handler(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("handling %q: %v", r.RequestURI, err)
			}
		} else {
			log.Println("User not logged in")

			http.Redirect(w, r, "/login/", 302)
		}

	})
}
