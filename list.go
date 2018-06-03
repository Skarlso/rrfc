package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const (
	// ListLocation is the location of the full list of RFCs.
	ListLocation = "https://www.ietf.org/download/rfc-index.txt"
	// FilePath is the path to the list file.
	FilePath = "list"
	// FileName is the name of the file the full RFC list is stored in.
	// This file is periodically update, once per day.
	FileName = "list.txt"
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

func parseList(list string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	re := regexp.MustCompile("^(\\d+) (.*?\\.)")
	filepath := filepath.Join(pwd, FilePath, list)
	f, err := os.Open(filepath)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		var (
			n    string
			desc string
		)
		match := re.FindStringSubmatch(l)
		if len(match) > 2 {
			n = match[1]
			desc = match[2]
			insertRFC(n, desc)
		}
	}

	return nil
}

func getRandomNumber() int {
	return 0
}
