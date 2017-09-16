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

	"github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal/pkg/conn"
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/logging"
	"github.com/fokal/fokal/pkg/ratelimit"
	"github.com/fokal/fokal/pkg/routes"
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
	RedisPass          string
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

	log.Println(cfg.PostgresURL)

	AppState.Vision, AppState.Maps, _ = conn.DialGoogleServices(cfg.GoogleToken)
	AppState.DB = conn.DialPostgres(cfg.PostgresURL)
	AppState.RD = conn.DialRedis(cfg.RedisURL, cfg.RedisPass)
	AppState.Local = cfg.Local
	AppState.Port = cfg.Port
	AppState.DB.SetMaxOpenConns(20)
	AppState.DB.SetMaxIdleConns(50)
	AppState.KeyHash = "554b5db484856bfa16e7da70a427dc4d9989678a"

	// RSA Keys
	AppState.PrivateKey, AppState.PublicKeys = ParseKeys()
	AppState.SessionLifetime = time.Hour * 16

	AppState.RefreshAt = time.Minute * 15

	// Refreshing Materialized View
	refreshMaterializedView()
	refreshGoogleOauthKeys()

	var secureMiddleware = secure.New(secure.Options{
		AllowedHosts:          []string{"api.fok.al", "dev.fok.al", "fok.al"},
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
		AllowedOrigins:     []string{"https://fok.al", "https://dev.fok.al", "http://localhost:3000"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		AllowedHeaders:     []string{"Authorization", "Content-Type"},
		AllowedMethods:     []string{"GET", "PUT", "OPTIONS", "PATCH", "POST"},
	})

	var base = alice.New(
		handler.SentryRecovery,
		ratelimit.RateLimit,
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
	api.NotFoundHandler = base.Then(http.HandlerFunc(handler.NotFound))

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

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if err != nil {
		log.Fatal(err)
	}
	parsedKeys[AppState.KeyHash] = publicKey

	privateStr := os.Getenv("PRIVATE_KEY")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateStr))
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, parsedKeys

}

func refreshMaterializedView() {
	tick := time.NewTicker(time.Minute * 15)
	go func() {
		for {
			select {
			case <-tick.C:
				AppState.DB.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY searches;")
			}
		}
	}()
}

func refreshGoogleOauthKeys() {
	tick := time.NewTicker(time.Minute * 10)
	go func() {
		for {
			select {
			case <-tick.C:
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
				for kid, pem := range keys {
					publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
					if err != nil {
						log.Fatal(err)
					}
					AppState.PublicKeys[kid] = publicKey
				}

			}
		}
	}()
}
