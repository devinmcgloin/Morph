package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fokal/fokal-core/pkg/conn"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {
	exampleDBQueryMultipleResultSets()
	postgresURL := os.Getenv("DATABASE_URL")
	db = conn.DialPostgres(postgresURL)
}

func exampleDBQueryMultipleResultSets() {
	rows, err := db.Query(
		`
	  select
		  id, shortcode from content.images
	  ;
	  select
		  id, 2
		  from content.images
	  ;
		  `)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id %d name is %s\n", id, name)
	}
	if !rows.NextResultSet() {
		log.Fatal("expected more result sets", rows.Err())
	}
	var roleMap = map[int64]string{
		1: "user",
		2: "admin",
		3: "gopher",
	}
	for rows.Next() {
		var (
			id   int64
			role int64
		)
		if err := rows.Scan(&id, &role); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id %d has role %s\n", id, roleMap[role])
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
