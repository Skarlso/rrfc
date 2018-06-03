package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// RFC contains an RFC retrieved from the database.
type RFC struct {
	Number      string
	Description string
}

func insertRFC(n string, desc string) error {
	connStr := "user=rrfc dbname=rrfc sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into rfcs (name, desc) values (?, ?)", n, desc)
	return err
}

func getRandomRow() (RFC, error) {
	connStr := "user=rrfc dbname=rrfc sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return RFC{}, err
	}
	rfc := RFC{}
	row, err := db.Query("select number, desc from rfcs order by random() limit 1")
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
