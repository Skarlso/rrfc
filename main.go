package main

import (
	"io/ioutil"
	"log"
)

func writeOutRandomRFC() {
	rfc, err := getRandomRow()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(".rfc", []byte(rfc.Number+":"+rfc.Description), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := createDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = downloadRFCList()
	if err != nil {
		log.Fatal(err)
	}
	parseListConcurrent("list.txt")
	writeOutRandomRFC()
}
