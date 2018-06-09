package main

import "log"

func main() {
	err := createDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = downloadRFCList()
	if err != nil {
		log.Fatal(err)
	}
	wipeRfcs()
	parseListConcurrent("list.txt")
	writeOutRandomRFC()
	writeOutAllPreviousRFCHTML()
}
