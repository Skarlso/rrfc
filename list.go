package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	// ListLocation is the location of the full list of RFCs.
	ListLocation = "https://www.ietf.org/download/rfc-index.txt"
	// FilePath is the path to the list file.
	FilePath = "list"
	// FileName is the name of the file the full RFC list is stored in.
	// This file is periodically update, once per day.
	FileName = "list.txt"
	// ChunkCount is the number of items the list is diveded into.
	ChunkCount = 100
)

func downloadRFCList() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filepath := filepath.Join(pwd, FilePath, FileName)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(ListLocation)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func parseListConcurrent(list string) error {
	var wg sync.WaitGroup
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	filepath := filepath.Join(pwd, FilePath, list)
	content, err := ioutil.ReadFile(filepath)
	split := strings.Split(string(content), "\n")
	tx, _ := beginTransaction()
	for {
		segment := len(split) / ChunkCount
		if segment == 0 {
			wg.Add(1)
			go handleSegment(split, tx, &wg)
			break
		}
		chunk := make([]string, 0)
		chunk, split = append(chunk, split[:segment]...), split[segment:]
		wg.Add(1)
		go handleSegment(chunk, tx, &wg)
	}
	wg.Wait()
	tx.Commit()
	return nil
}

func handleSegment(list []string, tx *sql.Tx, wg *sync.WaitGroup) {
	re := regexp.MustCompile("^(\\d+) (.*)")
	stmt, err := prepareStatement(tx)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i, s := range list {
		var (
			n    string
			desc string
		)
		match := re.FindStringSubmatch(s)
		if len(match) > 2 {
			n = match[1]
			desc = match[2]
			if !strings.ContainsAny(desc, ".") && i+1 < len(list) {
				desc = strings.Trim(desc, "\n")
				plus := strings.Trim(list[i+1], "\n")
				desc += " " + plus
			}
			// TODO: Needs to handle duplicate keys
			err = execStatement(stmt, n, desc)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	wg.Done()
}

func getRandomNumber() int {
	return 0
}
