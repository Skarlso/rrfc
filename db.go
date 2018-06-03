package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
)

// RFC contains an RFC retrieved from the database.
type RFC struct {
	Number      string
	Description string
}

var db *sql.DB

func init() {
	connStr := "user=rrfc dbname=rfcs sslmode=disable password=password123"
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("couldn't open db connection: ", err)
	}
	db = dbConn
}

func createDatabase() error {
	create := `
	CREATE TABLE rfcs (
		id SERIAL NOT NULL,
        number character varying(100) NOT NULL,
        description character varying(500) NOT NULL
	)
	`
	_, err := db.Exec(create)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Name() == "already exists" {
			fmt.Println("table already exists.")
			return nil
		}
	}

	return err
}

func insertRFC(n string, desc string) error {
	_, err := db.Exec("insert into rfcs (name, description) values (?, ?)", n, desc)
	return err
}

func getRandomRow() (RFC, error) {
	rfc := RFC{}
	row, err := db.Query("select number, description from rfcs order by random() limit 1")
	for row.Next() {
		var n string
		var desc string
		if err = row.Scan(&n, &desc); err != nil {
			return RFC{}, err
		}
		rfc.Number = n
		rfc.Description = desc
	}
	return rfc, nil
}
