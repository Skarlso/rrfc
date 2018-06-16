package main

import (
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

	"github.com/alecthomas/template"
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

// RFC contains an RFC retrieved from the database.
type RFC struct {
	Number      string
	Description string
	RFCList     []string
	Storage     Store
}

type rfcEntity struct {
	Number      string
	Description string
}

func (r *RFC) SetStore(store Store) {
	r.Storage = store
}

func (r *RFC) DownloadRFCList() error {
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

// ParseListConcurrent returns a list of parsed rfcEntities
// which will then be used as by a concurrent transaction save
// by the db.
func (r *RFC) parseListConcurrent(list string) []rfcEntity {
	pwd, err := os.Getwd()
	if err != nil {
		return []rfcEntity{}
	}
	filepath := filepath.Join(pwd, FilePath, list)
	content, err := ioutil.ReadFile(filepath)
	split := strings.Split(string(content), "\n")
	segmentChannels := make([]<-chan rfcEntity, 0)
	for {
		segment := len(split) / ChunkCount
		if segment == 0 {
			ch := handleSegment(split)
			segmentChannels = append(segmentChannels, ch)
			break
		}
		chunk := make([]string, 0)
		chunk, split = append(chunk, split[:segment]...), split[segment:]
		cs := handleSegment(chunk)
		segmentChannels = append(segmentChannels, cs)
	}

	rfcs := make([]rfcEntity, 0)
	rfcChannel := merge(segmentChannels...)
	for rfc := range rfcChannel {
		rfcs = append(rfcs, rfc)
	}
	return rfcs
}

func merge(in ...<-chan rfcEntity) <-chan rfcEntity {
	var wg sync.WaitGroup
	out := make(chan rfcEntity)
	wg.Add(len(in))
	output := func(ch <-chan rfcEntity) {
		for r := range ch {
			out <- r
		}
		wg.Done()
	}
	for _, c := range in {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func handleSegment(list []string) <-chan rfcEntity {
	retChannel := make(chan rfcEntity)
	re := regexp.MustCompile("^(\\d+) (.*)")
	go func() {
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
				r := rfcEntity{
					Number:      n,
					Description: desc,
				}
				retChannel <- r
			}
		}
		close(retChannel)
	}()
	return retChannel
}

func (r *RFC) WriteOutRandomRFC() {
	rfc, err := r.Storage.LoadRandom()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(".rfc", []byte(rfc.Number+":"+rfc.Description), 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = r.Storage.StoreRFC(rfc.Number, rfc.Description)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RFC) WriteOutAllPreviousRFCHTML() {
	rfcs, err := r.Storage.LoadAllPrevious()
	if err != nil {
		log.Fatal(err)
	}
	for _, rfc := range rfcs {
		filePath := filepath.Join("files", rfc.Number+".html")
		if _, osErr := os.Stat(filePath); osErr == nil {
			fmt.Println("skipping existing file: ", filePath)
			continue
		}
		rfcTemplate, _ := ioutil.ReadFile("rfc.template")
		t := template.Must(template.New("rfc").Parse(string(rfcTemplate)))
		f, err := os.Create(filePath)
		if err != nil {
			log.Fatal("error while creating file: ", err)
		}
		defer f.Close()
		err = t.Execute(f, rfc)
		if err != nil {
			log.Fatal("error writing file: ", err)
		}
	}
}
