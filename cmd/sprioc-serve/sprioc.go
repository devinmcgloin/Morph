package main

import (
	"flag"
	"log"
	"os"

	"github.com/sprioc/composer/pkg/daemon"
)

func ProcessFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.Host, "host-name", "localhost", "Host name to serve at")
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

	cfg.GoogleToken = googleToken
	cfg.PostgresURL = postgresURL

	daemon.Run(cfg)
}
