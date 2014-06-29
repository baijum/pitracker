package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres",
		`user=baiju
                 dbname=pitrackerlocal
                 sslmode=disable
                 password='passwd'`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Query(`INSERT INTO "member" (
            username, password) VALUES ($1, $2)`, "guest", "guest")

	if err != nil {
		log.Fatal(err)
	}
}
