package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
)

type PostgresStore struct {
}

func setup() {
	user := os.Getenv("PG_USER")
	dbName := os.Getenv("PG_DBNAME")
	password := os.Getenv("PG_PASSWORD")
	sslMode := os.Getenv("PG_SSLMODE")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", host, port, user, dbName, sslMode, password)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("couldn't open db connection: ", err)
	}
	db = dbConn
}

func createDatabase() error {
	createRfcs := `
	CREATE TABLE rfcs (
		id SERIAL NOT NULL,
        number character varying(100) NOT NULL UNIQUE,
		description character varying(500) NOT NULL
	)
	`
	createPreviousRfcs := `
	CREATE TABLE previous_rfcs (
		id SERIAL NOT NULL,
        number character varying(100) NOT NULL UNIQUE,
		description character varying(500) NOT NULL
	)
	`
	_, err := db.Exec(createRfcs)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Name() == "duplicate_table" {
			fmt.Println("table already exists.")
			fmt.Println("clearing")
			err = wipeRfcs()
		} else {
			return err
		}
	}

	_, err = db.Exec(createPreviousRfcs)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Name() == "duplicate_table" {
			fmt.Println("table already exists.")
			err = nil
		} else {
			return err
		}
	}

	return err
}

func prepareStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare("insert into rfcs (number, description) values ($1, $2)")
}

func beginTransaction() (*sql.Tx, error) {
	return db.Begin()
}

func storeRFC(n, desc string) error {
	// needs to handle duplicate keys.
	_, err := db.Exec("insert into previous_rfcs (number, description) values ($1, $2)", n, desc)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Name() == "unique_violation" {
			fmt.Println("skipping duplicate entry")
			return nil
		}
	}
	return err
}

func wipeRfcs() error {
	_, err := db.Exec("delete from rfcs")
	return err
}

func execStatement(stmt *sql.Stmt, n, desc string) error {
	// handle duplicate entries?
	_, err := stmt.Exec(n, desc)
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

func getAllPreviousRFCS() ([]RFC, error) {
	rfcs := make([]RFC, 0)
	row, err := db.Query("select number, description from previous_rfcs order by number asc")
	for row.Next() {
		var n string
		var desc string
		if err = row.Scan(&n, &desc); err != nil {
			return rfcs, err
		}
		rfc := RFC{}
		rfc.Number = n
		rfc.Description = desc
		rfcs = append(rfcs, rfc)
	}
	return rfcs, nil
}
