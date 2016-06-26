package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/sprioc/sprioc-core/pkg/authentication"
	h "github.com/sprioc/sprioc-core/pkg/handlers"
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
	api := router.PathPrefix("/v0").Subrouter()
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

func NotImplemented(w http.ResponseWriter, r *http.Request) h.Response {
	log.Printf("Not implemented called from %s", r.URL)
	return h.Response{Code: http.StatusNotImplemented, Message: "This endpoint is not implemented. It'll be here soon!"}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/index.html")
}

func secure(f func(http.ResponseWriter, *http.Request) h.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := authentication.CheckUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		context.Set(r, "auth", user)

		w.Header().Set("Content-Type", "application/json")

		resp := f(w, r)
		w.WriteHeader(resp.Code)
		if len(resp.Data) != 0 {
			w.Write(resp.Data)
		} else {
			w.Write(resp.Format())
		}
	}
}

func unsecure(f func(http.ResponseWriter, *http.Request) h.Response) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ip, port, _ := net.SplitHostPort(r.RemoteAddr)
		log.Println(ip, port)

		w.Header().Set("Content-Type", "application/json")

		resp := f(w, r)
		w.WriteHeader(resp.Code)
		if len(resp.Data) != 0 {
			w.Write(resp.Data)
		} else {
			w.Write(resp.Format())
		}
	}
}
