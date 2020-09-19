package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("v0.01")

	sp := NewSubmissionProcessor()
	rp := NewRequestProcessor(sp)

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", wwwHomePage)
	myRouter.HandleFunc("/submit", wwwSubmitForm)
	myRouter.HandleFunc("/api/submit", rp.apiSubmitSolutionHandler).Methods("POST")
	myRouter.HandleFunc("/api/submission/{problemName}/{id}", rp.apiReadSingleSubmission)

	go sp.Process()
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
