package main

import "log"

func main() {
	// err := createDatabase()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = downloadRFCList()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// parseListConcurrent("list.txt")
	// writeOutRandomRFC()
	// writeOutAllPreviousRFCHTML()
	rfc := new(RFC)
	rfcs := rfc.parseListConcurrent("list.txt")
	log.Println("len: ", len(rfcs))
}
