package daemon

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fokal/fokal-core/pkg/conn"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/middleware"
	raven "github.com/getsentry/raven-go"
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
	GoogleToken        string
	AWSAccessKeyId     string
	AWSSecretAccessKey string

	SentryURL string
}

var AppState handler.State

const PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsW3uHvJvqaaMIW8wKP2E
NI3oVRghsNwUV4VN+5UH2oMAEaYaHiUfOvhXXRjPZo3q8f+v3rS4R7gfJXe8efP0
3x87DRB1uJlNNS777xDISnTLzVAOFFkLOTL9bOTJBlb69yCRhHV1NdUIPCGWntWC
WdKZBJ2zHOQUQgPpAn31imsYlvmlrLEoGNqKOPUQjwdtxEqEYpZyN84Hj5/NIhTC
F6rU8FhReQzEL27BHPfbUwTWUApmtfvCtrSc9pVM3MtlsMOf4OfoGg65kF5HJ/S8
tKRtL24z48ya+ntjbwbE3A5pEswm/Vm19wd77qbY5UILLmNf0xMQfwrkT/IcnBoD
pQIDAQAB
-----END PUBLIC KEY-----`

func Run(cfg *Config) {
	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(flag)

	router := mux.NewRouter()
	api := router.PathPrefix("/v0/").Subrouter()

	log.Printf("Serving at http://%s:%d", cfg.Host, cfg.Port)
	err := raven.SetDSN(cfg.SentryURL)
	if err != nil {
		log.Fatal("Sentry IO not configured")
	}

	if cfg.Local {
		cfg.PostgresURL = cfg.PostgresURL + "?sslmode=disable"
	}

	Vision, Maps, _ := conn.DialGoogleServices(cfg.GoogleToken)
	DB := conn.DialPostgres(cfg.PostgresURL)
	RD := conn.DialRedis(cfg.RedisURL)

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(50)
	KeyHash := "554b5db484856bfa16e7da70a427dc4d9989678a"

	// RSA Keys
	SessionLifetime := time.Hour * 16

	RefreshAt := time.Minute * 15

	AppState.Local = cfg.Local
	AppState.Port = cfg.Port
	// Refreshing Materialized View

	var secureMiddleware = secure.New(secure.Options{
		AllowedHosts:          []string{"api.fok.al", "alpha.fok.al", "beta.fok.al", "fok.al"},
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLRedirect:           true,
		SSLHost:               "api.fok.al",
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
		AllowedOrigins:     []string{"https://fok.al", "https://beta.fok.al", "https://alpha.fok.al", "http://localhost:3000"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		AllowedHeaders:     []string{"Authorization", "Content-Type"},
		AllowedMethods:     []string{"GET", "PUT", "OPTIONS", "PATCH", "POST", "DELETE"},
	})

	var base = alice.New(
		middleware.SentryRecovery,
		middleware.RateLimit,
		crs.Handler,
		middleware.Timeout,
		middleware.IP, middleware.UUID, secureMiddleware.Handler,
		context.ClearHandler, handlers.CompressHandler, middleware.ContentTypeJSON)

	//  ROUTES

	api.NotFoundHandler = base.Then(http.HandlerFunc(handler.NotFound))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port),
		handlers.LoggingHandler(os.Stdout, router)))
}
