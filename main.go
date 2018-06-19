package main

import (
	"log"
	"path/filepath"
)

func main() {
	rfc := new(RFC)
	pgStore := new(PostgresStore)
	pgStore.Connect()
	rfc.SetStore(pgStore)

	err := pgStore.CreateStore()
	if err != nil {
		log.Fatal("error creating store:", err)
	}

	rfc.DownloadRFCList()
	rfcs := rfc.parseListConcurrent(filepath.Join(FilePath, "list.txt"))
	pgStore.StoreList(rfcs)

	rfc.WriteOutRandomRFC()
	rfc.WriteOutAllPreviousRFCHTML()
}
