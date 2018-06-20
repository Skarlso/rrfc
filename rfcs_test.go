package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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

func TestWriteOutRandomRFC(t *testing.T) {
	rfc := new(RFC)
	ds := new(dummyStore)
	ds.RFC = RFC{Number: "0001", Description: "Description"}
	ds.Error = nil
	rfc.SetStore(ds)
	rfc.WriteOutRandomRFC()
	content, err := ioutil.ReadFile(".rfc")
	if err != nil {
		t.Fatal(".rfc file not found")
	}
	if string(content) != "0001:Description" {
		t.Fatal("content did not equal expected. was: ", string(content))
	}
}

func TestWriteOutRandomRFCFailLoadRandom(t *testing.T) {
	called := false
	logFatal = func(...interface{}) {
		called = true
	}
	rfc := new(RFC)
	ds := new(dummyStore)
	ds.RFC = RFC{Number: "0001", Description: "Description"}
	ds.Error = errors.New("failed")
	rfc.SetStore(ds)
	rfc.WriteOutRandomRFC()
	if !called {
		t.Fatal("logFatal was not called")
	}
}

func TestWriteOutRandomRFCFailWritingFile(t *testing.T) {
	os.Setenv("RFC_FILENAME", "")
	called := false
	logFatal = func(...interface{}) {
		called = true
	}
	rfc := new(RFC)
	ds := new(dummyStore)
	ds.Error = nil
	rfc.SetStore(ds)
	rfc.WriteOutRandomRFC()
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
	rfc.WriteOutRandomRFC()
	if called {
		t.Fatal("logFatal was called")
	}
}
