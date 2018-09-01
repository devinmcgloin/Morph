package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/fokal/fokal-core/pkg/conn"
	"github.com/fokal/fokal-core/pkg/services/color"
	"github.com/jmoiron/sqlx"
)

type clr struct {
	Name string `toml:"name"`
	Hex  string `toml:"hex"`
}

type colors struct {
	Color []clr
}

func main() {
	var contentType, path string

	postgresURL := os.Getenv("DATABASE_URL")
	if postgresURL == "" {
		fmt.Fprintf(os.Stderr, "Postgres URL not set at POSTGRES_URL")
		os.Exit(1)
	}
	db := conn.DialPostgres(postgresURL)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-file>\n", filepath.Base(os.Args[0]))
	}
	flag.StringVar(&contentType, "type", "", "Type of colors loaded (shade|specific)")
	flag.StringVar(&path, "path", "", "Path to load")

	flag.Parse()

	if path == "" || contentType == "" {
		flag.Usage()
		os.Exit(1)
	}

	var content uint8
	if contentType == "specific" {
		content = color.SpecificColor
	} else if contentType == "shade" {
		content = color.Shade
	} else {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(path, contentType)
	if err := run(db, path, content); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

}

func run(db *sqlx.DB, file string, t uint8) error {
	var c colors
	toAdd := make(map[string]string)

	blob, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if _, err := toml.Decode(string(blob), &c); err != nil {
		return err
	}

	fmt.Printf("Path %s contains %d colors.\n", file, len(c.Color))
	for _, clr := range c.Color {
		toAdd[clr.Hex] = clr.Name
	}

	table := color.NewWithType(db, t)

	err = table.AddColors(toAdd)
	return err
}
