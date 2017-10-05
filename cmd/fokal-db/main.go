package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/fokal/fokal/pkg/conn"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	postgresURL := os.Getenv("DATABASE_URL")
	db = conn.DialPostgres(postgresURL)
}

func main() {
	ExampleDB_Query_multipleResultSets()

}
func ExampleDB_Query() {
	age := 27
	rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func ExampleDB_QueryRow() {
	id := 123
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Username is %s\n", username)
	}
}

func ExampleDB_Query_multipleResultSets() {
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
