package daemon

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/sprioc/composer/pkg/cache"
	"github.com/sprioc/composer/pkg/metadata"
	"github.com/sprioc/composer/pkg/rsp"
	"github.com/sprioc/composer/pkg/sql"
)

type Config struct {
	Port               int
	Host               string
	PostgresURL        string
	RedisURL           string
	RedisPass          string
	GoogleToken        string
	AWSAccessKeyId     string
	AWSSecretAccessKey string
}

func Run(cfg *Config) {
	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flag)

	router := mux.NewRouter()
	api := router.PathPrefix("/v0/").Subrouter()

	log.Printf("Serving at http://%s:%d", cfg.Host, cfg.Port)

	//  ROUTES
	registerImageRoutes(api)
	registerUserRoutes(api)
	// registerCollectionRoutes(api)
	// registerSearchRoutes(api)
	// registerLuckyRoutes(api)
	registerAuthRoutes(api)

	metadata.Configure(cfg.GoogleToken)
	sql.Configure(cfg.PostgresURL)
	cache.Configure(cfg.RedisURL, cfg.RedisPass)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), context.ClearHandler(handlers.LoggingHandler(os.Stdout,
		handlers.CompressHandler(router)))))
}

// NotImplemented returns the standard response for endpoints that have not been implemented
func NotImplemented(w http.ResponseWriter, r *http.Request) rsp.Response {
	log.Printf("Not implemented called from %s", r.URL)
	return rsp.Response{Code: http.StatusNotImplemented, Message: "This endpoint is not implemented. It'll be here soon!"}
}
