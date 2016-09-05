package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/sprioc/composer/pkg/rsp"
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
	api := router.PathPrefix("/v0/").Subrouter()
	port := os.Getenv("PORT")

	log.Printf("Serving at http://localhost:%s", port)

	//  ROUTES
	registerImageRoutes(api)
	// registerUserRoutes(api)
	// registerCollectionRoutes(api)
	// registerSearchRoutes(api)
	// registerLuckyRoutes(api)
	// registerAuthRoutes(api)

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout,
		handlers.CompressHandler(router))))
}

// NotImplemented returns the standard response for endpoints that have not been implemented
func NotImplemented(w http.ResponseWriter, r *http.Request) rsp.Response {
	log.Printf("Not implemented called from %s", r.URL)
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This endpoint is not implemented. It'll be here soon!"}
}
