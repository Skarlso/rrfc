package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {

}

func GetRandomRFC(w http.ResponseWriter, r *http.Request) {

	// res := struct {
	// 	Number int
	// }
}

func serve() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/r", GetRandomRFC).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
