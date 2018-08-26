package daemon

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fokal/fokal-core/pkg/services/stream"
	"github.com/fokal/fokal-core/pkg/services/tag"
	"github.com/fokal/fokal-core/pkg/services/user"
	"github.com/fokal/fokal-core/pkg/services/vision"

	"github.com/fokal/fokal-core/pkg/services/permission"

	"github.com/fokal/fokal-core/pkg/conn"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/middleware"
	"github.com/fokal/fokal-core/pkg/services/authentication"
	"github.com/fokal/fokal-core/pkg/services/cache"
	"github.com/fokal/fokal-core/pkg/services/color"
	"github.com/fokal/fokal-core/pkg/services/upload"
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

	VisionService, MapService, _ := conn.DialGoogleServices(cfg.GoogleToken)
	DB := conn.DialPostgres(cfg.PostgresURL)
	RD := conn.DialRedis(cfg.RedisURL)

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(50)

	// RSA Keys
	SessionLifetime := time.Hour * 16

	CacheExpiryTime := time.Hour * 2

	AppState.Local = cfg.Local
	AppState.Port = cfg.Port

	AppState.CacheService = cache.New(RD, "cache:", CacheExpiryTime)
	AppState.ColorService = color.New(DB)
	AppState.StorageService = handler.StorageState{
		Content: upload.New("fokal-content", "us-west-1", "content"),
		Avatar:  upload.New("fokal-content", "us-west-1", "avatar"),
	}
	AppState.PermissionService = permission.New(DB)
	AppState.TagService = tag.New(DB)
	AppState.VisionService = vision.New(DB, VisionService)
	// AppState.SearchService = search.New(DB)

	fmt.Println(MapService)
	AppState.UserService = user.New(DB, AppState.PermissionService, AppState.ImageService)
	AppState.AuthService = authentication.New(DB, AppState.UserService, SessionLifetime)
	AppState.StreamService = stream.New(DB, AppState.ImageService, AppState.PermissionService)

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
