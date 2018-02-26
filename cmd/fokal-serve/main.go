package main

import (
	"flag"
	"log"
	"os"

	"github.com/fokal/fokal/pkg/daemon"

	"strconv"
)

func ProcessFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.Host, "host", "localhost", "Host name to serve at")
	flag.IntVar(&cfg.Port, "port", 8080, "Port to Listen on")
	flag.BoolVar(&cfg.Local, "local", false, "True if running locally")

	flag.Parse()
	return cfg
}

func main() {
	cfg := ProcessFlags()

	port := os.Getenv("PORT")
	if port != "" {
		p, _ := strconv.ParseInt(port, 10, 32)
		cfg.Port = int(p)
	}

	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		log.Fatal("Postgres URL not set at DATABASE_URL")
	}

	googleToken := os.Getenv("GOOGLE_API_TOKEN")
	if googleToken == "" {
		log.Fatal("Google API Token not set at GOOGLE_API_TOKEN")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("Redis URL not set at REDIS_URL")
	}

	redisPass := os.Getenv("REDIS_PASS")
	if redisPass == "" {
		log.Fatal("Redis Pass not set at REDIS_PASS")
	}

	// AWS auth
	AWSAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if AWSAccessKey == "" {
		log.Fatal("AWS Access Key Id not set at AWS_ACCESS_KEY_ID")
	}

	AWSSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if AWSSecret == "" {
		log.Fatal("AWS Secret Access Key not set at AWS_SECRET_ACCESS_KEY")
	}

	SentryURL := os.Getenv("SENTRY_URL")
	if SentryURL == "" {
		log.Fatal("Sentry URL not set at SENTRY_URL")
	}

	cfg.GoogleToken = googleToken
	cfg.PostgresURL = postgresURL
	cfg.RedisURL = redisURL
	cfg.RedisPass = redisPass
	cfg.AWSAccessKeyId = AWSAccessKey
	cfg.AWSSecretAccessKey = AWSSecret
	cfg.SentryURL = SentryURL

	daemon.Run(cfg)
}
