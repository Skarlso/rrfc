package main

import "log"

func main() {
	err := createDatabase()
	if err != nil {
		log.Fatal(err)
	}
	parseList("sample.txt")
}
