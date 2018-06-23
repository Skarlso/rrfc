package main

import (
	"log"
	"reflect"
	"testing"
)

var skip = false
var psql *PostgresStore

func init() {
	psql = new(PostgresStore)
	psql.Connect()
	err := psql.Ping()
	if err != nil {
		skip = true
	}
}

func cleanup() {
	log.Println("running cleanup")
	_, err := psql.Exec("drop table if exists rfcs")
	if err != nil {
		log.Println("filed to cleanup: ", err)
	}
	_, err = psql.Exec("drop table if exists previous_rfcs")
	if err != nil {
		log.Println("filed to cleanup: ", err)
	}
}

func TestPostgresCreateStore(t *testing.T) {
	if skip {
		t.Skip("skipping because postgres is not runnig")
	}
	defer cleanup()
	err := psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
}

func TestPostgresCreateStoreIsIdenpotent(t *testing.T) {
	if skip {
		t.Skip("skipping because postgres is not runnig")
	}
	defer cleanup()
	err := psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	err = psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
}

func TestInsertPreviousRFC(t *testing.T) {
	if skip {
		t.Skip("skipping because postgres is not runnig")
	}
	defer cleanup()
	err := psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	err = psql.StorePreviousRFC("0001", "Description")
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	row, err := psql.Query("select number, description from previous_rfcs")
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	rfc := new(RFC)
	for row.Next() {
		var n string
		var desc string
		if err = row.Scan(&n, &desc); err != nil {
			t.Fatal("error reading row: ", err)
		}
		rfc.Number = n
		rfc.Description = desc
	}
	if rfc.Number != "0001" || rfc.Description != "Description" {
		t.Fatalf("loaded rfc did not match expected. was: %+v\n", rfc)
	}
}

func TestLoadRFCListAndLoadRandom(t *testing.T) {
	if skip {
		t.Skip("skipping because postgres is not runnig")
	}
	defer cleanup()
	err := psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	rfcs := []rfcEntity{
		rfcEntity{Number: "0001", Description: "description"},
	}
	psql.StoreList(rfcs)
	rfc, err := psql.LoadRandom()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	if rfc.Number != "0001" || rfc.Description != "description" {
		t.Fatalf("loaded rfc did not match expected. was: %+v\n", rfc)
	}
}

func TestLoadRandomLoadsADifferentRFC(t *testing.T) {
	if skip {
		t.Skip("skipping because postgres is not runnig")
	}
	defer cleanup()
	err := psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	rfcs := []rfcEntity{
		rfcEntity{Number: "0001", Description: "description1"},
		rfcEntity{Number: "0002", Description: "description2"},
		rfcEntity{Number: "0003", Description: "description3"},
		rfcEntity{Number: "0004", Description: "description4"},
	}
	psql.StoreList(rfcs)
	notRandom := true
	rfc, _ := psql.LoadRandom()
	for i := 0; i < 100; i++ {
		currRfc, _ := psql.LoadRandom()
		if !reflect.DeepEqual(rfc, currRfc) {
			notRandom = false
			break
		}
		rfc = currRfc
	}
	if notRandom {
		t.Fatal("the loading of random rfcs was not random after a 100 tries")
	}
}

func TestLoadAllPreviousRFCS(t *testing.T) {
	if skip {
		t.Skip("skipping because postgres is not runnig")
	}
	defer cleanup()
	err := psql.CreateStore()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	err = psql.StorePreviousRFC("0001", "Description")
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	rfcs, err := psql.LoadAllPrevious()
	if err != nil {
		t.Fatal("did not expect error: ", err)
	}
	expected := []RFC{
		RFC{Number: "0001", Description: "Description"},
	}
	if !reflect.DeepEqual(expected, rfcs) {
		t.Fatal("loaded list did not equal expected. was:", rfcs)
	}
}
