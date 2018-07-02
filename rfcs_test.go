package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type dummyStore struct {
	Error error
	RFC   RFC
	RFCS  []RFC
}

func (ds *dummyStore) StorePreviousRFC(string, string) error {
	return ds.Error
}
func (ds *dummyStore) StoreList([]rfcEntity) {

}
func (ds *dummyStore) LoadRandom() (RFC, error) {
	return ds.RFC, ds.Error
}
func (ds *dummyStore) LoadAllPrevious() ([]RFC, error) {
	return ds.RFCS, ds.Error
}
func (ds *dummyStore) Wipe() error {
	return ds.Error
}
func (ds *dummyStore) CreateStore() error {
	return ds.Error
}
func (ds *dummyStore) Connect() {

}

func TestParseListConcurrent(t *testing.T) {
	expected := []rfcEntity{
		rfcEntity{"0001", "Host Software. S. Crocker. April 1969. (Format: TXT=21088 bytes)"},
		rfcEntity{"0002", "Host software. B. Duvall. April 1969. (Format: TXT=17145 bytes)"},
		rfcEntity{"0003", "Documentation conventions. S.D. Crocker. April 1969. (Format:"},
		rfcEntity{"0004", "Network timetable. E.B. Shapiro. March 1969. (Format: TXT=5933"},
		rfcEntity{"0005", "Decode Encode Language (DEL). J. Rulifson. June 1969. (Format:"},
	}
	rfc := new(RFC)
	res := rfc.parseListConcurrent(filepath.Join("fixture", "test_parse_list_con.txt"))
	if !reflect.DeepEqual(expected, res) {
		t.Fatalf("expected \n'%+v'\n did not match actual \n'%+v'\n", expected, res)
	}
}

func TestGetRandomRFC(t *testing.T) {
	rfc := new(RFC)
	ds := new(dummyStore)
	ds.RFC = RFC{Number: "0001", Description: "Description"}
	ds.Error = nil
	rfc.SetStore(ds)
	rfc.GetRandomRFC()
	content, err := ioutil.ReadFile(".rfc")
	if err != nil {
		t.Fatal(".rfc file not found")
	}
	if string(content) != "0001:Description" {
		t.Fatal("content did not equal expected. was: ", string(content))
	}
}

func TestGetRandomRFCFailLoadRandom(t *testing.T) {
	called := false
	logFatal = func(...interface{}) {
		called = true
	}
	rfc := new(RFC)
	ds := new(dummyStore)
	ds.RFC = RFC{Number: "0001", Description: "Description"}
	ds.Error = errors.New("failed")
	rfc.SetStore(ds)
	rfc.GetRandomRFC()
	if !called {
		t.Fatal("logFatal was not called")
	}
}

func TestWriteOutRandomStoringRFCWorks(t *testing.T) {
	os.Setenv("RFC_FILENAME", ".rfc")
	called := false
	logFatal = func(...interface{}) {
		called = true
	}
	rfc := new(RFC)
	ds := new(dummyStore)
	ds.RFC = RFC{Number: "0001", Description: "Description"}
	ds.Error = nil
	rfc.SetStore(ds)
	rfc.GetRandomRFC()
	if called {
		t.Fatal("logFatal was called")
	}
}

func TestFileDownload(t *testing.T) {
	called := false
	logFatal = func(args ...interface{}) {
		called = true
		log.Println(args)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "test_content")
	}))
	defer ts.Close()
	os.Setenv("LIST_URL", ts.URL)
	rfc := new(RFC)
	rfc.DownloadRFCList()
	if called {
		t.Fatal("logFatal was called")
	}
}

func TestFileDownloadFailedHttpStatus(t *testing.T) {
	called := false
	var message string
	logFatal = func(args ...interface{}) {
		called = true
		message = args[0].(string)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "danger", 500)
	}))
	defer ts.Close()
	os.Setenv("LIST_URL", ts.URL)
	rfc := new(RFC)
	rfc.DownloadRFCList()
	if !called {
		t.Fatal("logFatal was not called")
	}
	if message != "bad status: 500 Internal Server Error" {
		t.Fatal("error not picked up")
	}
}

func cleanupFiles(file string) {
	os.Remove(filepath.Join("fixture", "files", file))
}

func TestWriteOutAllPreviousRFCHTML(t *testing.T) {
	os.Setenv("SITE_LOCATION", "fixture")
	defer cleanupFiles("0001.html")
	rfc := new(RFC)
	psql := new(dummyStore)
	rfcs := []RFC{
		RFC{Number: "0001", Description: "Description"},
	}
	psql.RFCS = rfcs
	rfc.SetStore(psql)
	rfc.WriteOutAllPreviousRFCHTML()
	if _, err := os.Stat(filepath.Join("fixture", "files", "0001.html")); err != nil {
		t.Fatal("generated file at expected location not found")
	}
}

func TestWriteOutAllPreviousRFCHTMLSkippingExistingFiles(t *testing.T) {
	os.Setenv("SITE_LOCATION", "fixture")
	defer cleanupFiles("0001.html")
	rfc := new(RFC)
	psql := new(dummyStore)
	rfcs := []RFC{
		RFC{Number: "0001", Description: "Description"},
	}
	psql.RFCS = rfcs
	rfc.SetStore(psql)
	rfc.WriteOutAllPreviousRFCHTML()
	if _, err := os.Stat(filepath.Join("fixture", "files", "0001.html")); err != nil {
		t.Fatal("generated file at expected location not found")
	}
	rfc.WriteOutAllPreviousRFCHTML()
}
