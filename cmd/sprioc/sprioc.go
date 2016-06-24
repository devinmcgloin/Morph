package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/devinmcgloin/sprioc/src/api/auth"
	"github.com/devinmcgloin/sprioc/src/spriocError"
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
	registerSearchRoutes(api)
	registerLuckyRoutes(api)
	registerAuthRoutes(router)

	router.HandleFunc("/", serveHTML)

	// ASSETS
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(router))))
}

func NotImplemented(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Not implemented called from %s", r.URL)
	http.Error(w, "Not Implemented", 509)
	return nil
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/index.html")
}

func secureWrappedMiddle(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := auth.CheckUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		context.Set(r, "auth", user)

		err = f(w, r)
		if err != nil {
			custErr := err.(spriocError.SpriocError)
			if custErr.Code != 0 {
				http.Error(w, custErr.Error(), custErr.Code)
				log.Printf("error handling %q: %v", r.RequestURI, err)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("unknown error handling %q: %v", r.RequestURI, err)
			return
		}
	}
}

func unsecureWrappedMiddle(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			custErr := err.(spriocError.SpriocError)
			if custErr.Code != 0 {
				http.Error(w, custErr.Error(), custErr.Code)
				log.Printf("error handling %q: %v", r.RequestURI, err)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("unknown error handling %q: %v", r.RequestURI, err)
			return
		}
	}
}
