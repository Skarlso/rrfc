package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
)

// PostgresStore is a Store backed by Postgres type database.
type PostgresStore struct {
	*sql.DB
}

// Connect to a store.
func (ps *PostgresStore) Connect() {
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
	ps.DB = dbConn
}

// CreateStore creates a store backend to be used. This function contains setup
// for any kind of store, ie. create tables or files.
func (ps *PostgresStore) CreateStore() error {
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
	_, err := ps.Exec(createRfcs)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Name() == "duplicate_table" {
			fmt.Println("table already exists.")
			fmt.Println("clearing")
			err = ps.Wipe()
		} else {
			return err
		}
	}

	_, err = ps.Exec(createPreviousRfcs)
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

// StoreRFC stores a single RFC.
func (ps *PostgresStore) StoreRFC(n, desc string) error {
	// needs to handle duplicate keys.
	_, err := ps.Exec("insert into previous_rfcs (number, description) values ($1, $2)", n, desc)
	if e, ok := err.(*pq.Error); ok {
		if e.Code.Name() == "unique_violation" {
			fmt.Println("skipping duplicate entry")
			return nil
		}
	}
	return err
}

// Wipe wipes the main database. The idea is rather than doing an upsert or a
// checked insert, we just clear everything and store everything on the weekly cycle.
func (ps *PostgresStore) Wipe() error {
	_, err := ps.Exec("delete from rfcs")
	return err
}

// LoadRandom gives back a random RFC.
func (ps *PostgresStore) LoadRandom() (RFC, error) {
	rfc := RFC{}
	row, err := ps.Query("select number, description from rfcs order by random() limit 1")
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

// StoreList stores an entire list of RFCs. This is done via a transaction here
// so that it is a single commit instead of several thousand.
func (ps *PostgresStore) StoreList(rfcs []rfcEntity) {
	tx, err := ps.Begin()
	if err != nil {
		log.Fatal("error beginning transaction: ", err)
	}
	stmt, err := tx.Prepare("insert into rfcs (number, description) values ($1, $2)")
	if err != nil {
		log.Fatal("error while preparing statement: ", err)
	}
	for _, r := range rfcs {
		_, err := stmt.Exec(r.Number, r.Description)
		if err != nil {
			log.Fatal("error while executing statement: ", err)
		}
	}
	tx.Commit()
}

// LoadAllPrevious is used to load all previous randomly selected rfcs in order
// to create a static html file for them.
func (ps *PostgresStore) LoadAllPrevious() ([]RFC, error) {
	rfcs := make([]RFC, 0)
	row, err := ps.Query("select number, description from previous_rfcs order by number asc")
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
