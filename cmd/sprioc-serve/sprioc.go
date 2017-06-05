package main

import (
	"flag"
	"log"
	"os"

	"github.com/sprioc/composer/pkg/daemon"
)

func ProcessFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.Host, "host", "localhost", "Host name to serve at")
	flag.IntVar(&cfg.Port, "port", 8080, "Port to Listen on")

	flag.Parse()
	return cfg
}
func main() {
	cfg := ProcessFlags()

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		log.Fatal("Postgres URL not set at POSTGRES_URL")
	}

	googleToken := os.Getenv("GOOGLE_API_TOKEN")
	if googleToken == "" {
		log.Fatal("Google API Token not set at GOOGLE_API_TOKEN")
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("Redis URL not set at REDIS_URL")
	}

	cfg.GoogleToken = googleToken
	cfg.PostgresURL = postgresURL
	cfg.RedisURL = redisURL

	daemon.Run(cfg)
}
