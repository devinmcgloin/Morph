package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devinmcgloin/sprioc/src/api/auth"
	"github.com/devinmcgloin/sprioc/src/api/session"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Port must be set")
	}

	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flag)

}

func main() {

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v0").Subrouter()
	port := os.Getenv("PORT")

	log.Printf("Serving at http://localhost:%s", port)

	//  ROUTES
	registerImageRoutes(api)
	registerUserRoutes(api)
	registerCollectionRoutes(api)
	registerAlbumRoutes(api)

	router.HandleFunc("/", serveHTML)

	// ASSETS
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)))
}

func middleCheckAuth(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
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

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not implemented called from %s", r.URL)
	http.Error(w, "Not Implemented", 404)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/index.html")
}
