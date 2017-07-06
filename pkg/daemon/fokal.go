package daemon

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/devinmcgloin/fokal/pkg/conn"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/logging"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
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

var AppState handler.State

func Run(cfg *Config) {
	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flag)

	router := mux.NewRouter()
	api := router.PathPrefix("/v0/").Subrouter()

	log.Printf("Serving at http://%s:%d", cfg.Host, cfg.Port)

	AppState.Vision, AppState.Maps, _ = conn.DialGoogleServices(cfg.GoogleToken)
	AppState.DB = conn.DialPostgres(cfg.PostgresURL)
	AppState.RD = conn.DialRedis(cfg.RedisURL, cfg.RedisPass)
	AppState.Local = cfg.Local

	var secureMiddleware = secure.New(secure.Options{
		AllowedHosts:          []string{"api.sprioc.xyz"},
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLRedirect:           true,
		SSLHost:               "api.sprioc.xyz",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IsDevelopment:         AppState.Local,
	})

	var crs = cors.New(cors.Options{
		AllowedOrigins:   []string{"https://sprioc.xyz"},
		AllowCredentials: true,
	})

	var base = alice.New(logging.IP, logging.UUID, secureMiddleware.Handler, crs.Handler,
		context.ClearHandler, handlers.CompressHandler, logging.ContentTypeJSON, handler.Timeout)

	//  ROUTES
	registerImageRoutes(api, base)
	//registerUserRoutes(api)
	// registerCollectionRoutes(api)
	// registerSearchRoutes(api)
	// registerLuckyRoutes(api)
	//registerAuthRoutes(api)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port),
		handlers.LoggingHandler(os.Stdout, router)))
}
