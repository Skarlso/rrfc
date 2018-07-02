package main

import (
	"log"
	"path/filepath"
)

var logFatal = log.Fatal

func main() {
	rfc := new(RFC)
	pgStore := new(PostgresStore)
	pgStore.Connect()
	rfc.SetStore(pgStore)

	err := pgStore.CreateStore()
	if err != nil {
		logFatal("error creating store:", err)
	}

	rfc.DownloadRFCList()
	rfcs := rfc.parseListConcurrent(filepath.Join(FilePath, "list.txt"))
	pgStore.StoreList(rfcs)

	rrfc := rfc.GetRandomRFC()
	rfc.WriteOutAllPreviousRFCHTML()
	rfc.WriteOutPreviousHTML()
	rfc.WriteOutIndexHTML(rrfc)
}
