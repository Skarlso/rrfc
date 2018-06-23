package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/alecthomas/template"
	_ "github.com/joho/godotenv/autoload"
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

// SetStore sets up a store implementation for RFC to use
// Use DUMMY store here in order to unit test RFC.
func (r *RFC) SetStore(store Store) {
	r.Storage = store
}

// DownloadRFCList gets a list of all available RFCs
func (r *RFC) DownloadRFCList() {
	listLocation := os.Getenv("LIST_LOCATION")
	pwd, err := os.Getwd()
	if err != nil {
		logFatal(err)
	}

	filepath := filepath.Join(pwd, FilePath, FileName)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		logFatal(err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(listLocation)
	if err != nil {
		logFatal(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf("bad status: %s", resp.Status)
		logFatal(message)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logFatal(err)
	}
}

// ParseListConcurrent returns a list of parsed rfcEntities
// which will then be used as by a concurrent transaction save
// by the db.
func (r *RFC) parseListConcurrent(list string) []rfcEntity {
	pwd, err := os.Getwd()
	if err != nil {
		return []rfcEntity{}
	}
	filepath := filepath.Join(pwd, list)
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

// WriteOutRandomRFC creates a file called .rfc where a random rfc is stored.
// This is done so the web application can get the number from a file rather then
// having to implement database access in the main semi-static HTML file.
func (r *RFC) WriteOutRandomRFC() {
	rfc, err := r.Storage.LoadRandom()
	if err != nil {
		logFatal(err)
	}
	rfcFilename := os.Getenv("RFC_FILENAME")
	err = ioutil.WriteFile(rfcFilename, []byte(rfc.Number+":"+rfc.Description), 0644)
	if err != nil {
		logFatal(err)
	}
	err = r.Storage.StorePreviousRFC(rfc.Number, rfc.Description)
	if err != nil {
		logFatal(err)
	}
}

// WriteOutAllPreviousRFCHTML creates HTML static files for all previous RFCs
// that were once selected. This is so disqus has a permanent link to point to
// when viewing past conversations and for convenience.
func (r *RFC) WriteOutAllPreviousRFCHTML() {
	rfcs, err := r.Storage.LoadAllPrevious()
	if err != nil {
		logFatal(err)
	}
	base := os.Getenv("SITE_LOCATION")
	for _, rfc := range rfcs {
		filePath := filepath.Join(base, "files", rfc.Number+".html")
		if _, osErr := os.Stat(filePath); osErr == nil {
			fmt.Println("skipping existing file: ", filePath)
			continue
		}
		rfcTemplate, _ := ioutil.ReadFile("rfc.template")
		t := template.Must(template.New("rfc").Parse(string(rfcTemplate)))
		f, err := os.Create(filePath)
		if err != nil {
			logFatal("error while creating file: ", err)
		}
		defer f.Close()
		err = t.Execute(f, rfc)
		if err != nil {
			logFatal("error writing file: ", err)
		}
	}
}
