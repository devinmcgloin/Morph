package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

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

	router := mux.NewRouter().Host("api.sprioc.xyz").Subrouter()
	api := router.PathPrefix("/v0").Subrouter()
	port := os.Getenv("PORT")

	log.Printf("Serving at http://localhost:%s", port)

	//  ROUTES
	registerImageRoutes(api)
	registerUserRoutes(api)
	registerCollectionRoutes(api)
	registerSearchRoutes(api)
	registerLuckyRoutes(api)
	registerAuthRoutes(api)

	router.HandleFunc("/", status)
	router.NotFoundHandler = http.HandlerFunc(notFound)

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(router))))
}

func NotImplemented(w http.ResponseWriter, r *http.Request) h.Response {
	log.Printf("Not implemented called from %s", r.URL)
	return h.Response{Code: http.StatusNotImplemented, Message: "This endpoint is not implemented. It'll be here soon!"}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	var response = make(map[string]interface{})
	response["message"] = "Not Found"
	response["code"] = 404
	response["documentation_url"] = "http://github.com/sprioc/sprioc-core"
	bytes, _ := json.MarshalIndent(response, "", "    ")
	w.Write(bytes)
	return
}

func status(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}
	m["status"] = "good"
	m["time"] = time.Now().Format(time.RFC3339)
	m["version"] = "v0"
	json, err := json.Marshal(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(json)
}
