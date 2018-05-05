package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fokal/fokal-core/pkg/conn"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: Provide DB URL under rel flag.")
	}

	var dbURL string

	flag.StringVar(&dbURL, "rel", "", "Database to Initialize")

	flag.Parse()

	if dbURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	db := conn.DialPostgres(dbURL)

	db.MustExec(createScript)
}
