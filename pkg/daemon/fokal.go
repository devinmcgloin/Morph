package daemon

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/devinmcgloin/fokal/pkg/conn"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	Port int
	Host string

	Local bool

	PostgresURL        string
	RedisURL           string
	RedisPass          string
	GoogleToken        string
	AWSAccessKeyId     string
	AWSSecretAccessKey string
}

var AppState State

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

	AppState.Vision, AppState.Maps, _ = conn.DialGoogleServices(cfg.GoogleToken)
	AppState.DB = conn.DialPostgres(cfg.PostgresURL)
	AppState.RD = conn.DialRedis(cfg.RedisURL, cfg.RedisPass)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), context.ClearHandler(handlers.LoggingHandler(os.Stdout,
		handlers.CompressHandler(router)))))
}
