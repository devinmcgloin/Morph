package daemon

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"crypto/rsa"
	"encoding/json"
	"io/ioutil"

	"time"

	"github.com/devinmcgloin/fokal/pkg/conn"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/logging"
	"github.com/devinmcgloin/fokal/pkg/routes"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/dgrijalva/jwt-go"
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
	AppState.Port = cfg.Port
	AppState.DB.SetMaxOpenConns(20)
	AppState.DB.SetMaxIdleConns(50)

	// RSA Keys
	AppState.PrivateKey, AppState.PublicKeys = ParseKeys()
	AppState.SessionLifetime = time.Hour * 16

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
		AllowedOrigins:     []string{"https://sprioc.xyz", "https://dev.sprioc.xyz", "http://localhost:3000"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		AllowedHeaders:     []string{"Authorization", "Content-Type"},
		AllowedMethods:     []string{"GET", "PUT", "OPTIONS", "PATCH", "POST"},
	})

	var base = alice.New(
		crs.Handler,
		handler.Timeout,
		logging.IP, logging.UUID, secureMiddleware.Handler,
		context.ClearHandler, handlers.CompressHandler, logging.ContentTypeJSON)

	//  ROUTES
	routes.RegisterCreateRoutes(&AppState, api, base)
	routes.RegisterModificationRoutes(&AppState, api, base)
	routes.RegisterRetrievalRoutes(&AppState, api, base)
	routes.RegisterSocialRoutes(&AppState, api, base)
	routes.RegisterSearchRoutes(&AppState, api, base)
	routes.RegisterRandomRoutes(&AppState, api, base)
	routes.RegisterAuthRoutes(&AppState, api, base)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port),
		handlers.LoggingHandler(os.Stdout, router)))
}

func ParseKeys() (*rsa.PrivateKey, map[string]*rsa.PublicKey) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	keys := make(map[string]string)
	err = json.Unmarshal(body, &keys)
	if err != nil {
		log.Fatal(err)
	}

	parsedKeys := make(map[string]*rsa.PublicKey)

	for kid, pem := range keys {
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			log.Fatal(err)
		}
		parsedKeys[kid] = publicKey
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(security.PublicKey))
	if err != nil {
		log.Fatal(err)
	}
	parsedKeys["554b5db484856bfa16e7da70a427dc4d9989678a"] = publicKey

	privateStr := os.Getenv("PRIVATE_KEY")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateStr))
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, parsedKeys

}
