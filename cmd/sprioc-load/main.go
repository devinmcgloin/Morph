package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/sprioc/composer/pkg/sql"
)

type color struct {
	Name string `toml:"name"`
	Hex  string `toml:"hex"`
}

type colors struct {
	Color []color
}

func main() {
	var contentType, path string

	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		fmt.Fprintf(os.Stderr, "Postgres URL not set at POSTGRES_URL")
		os.Exit(1)
	}
	sql.Configure(postgresURL)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-file>\n", filepath.Base(os.Args[0]))
	}
	flag.StringVar(&contentType, "type", "", "Type of colors loaded")
	flag.StringVar(&path, "path", "", "Path to load")

	flag.Parse()

	if path == "" || contentType == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(path, contentType)
	if err := run(path, contentType); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

}

func run(file, t string) error {
	var colors colors
	toAdd := make(map[string]string)

	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if _, err := toml.Decode(string(blob), &colors); err != nil {
		return err
	}

	fmt.Printf("Path %s contains %d colors.\n", file, len(colors.Color))
	for _, clr := range colors.Color {
		toAdd[clr.Hex] = clr.Name
	}
	err = sql.AddColors(toAdd, t)
	if err != nil {
		return err
	}
	return nil
}
