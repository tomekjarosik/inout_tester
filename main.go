package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomekjarosik/inout_tester/internal/submission"
)

// TODO: Add config for:
//   - problems (e.g. zad1, zad2 etc)
// TODO: Handle Ctrl+C properly
func main() {
	fmt.Println("Started new server at http://localhost:8080")

	storage := submission.NewDefaultStorage("submissions")
	sp := submission.NewProcessor(storage)
	rp := NewRequestProcessor(storage, sp)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", rp.RenderHomePage)
	myRouter.HandleFunc("/submit", rp.wwwSubmitForm)
	myRouter.HandleFunc("/api/submit", rp.apiSubmitSolutionHandler).Methods("POST")
	myRouter.HandleFunc("/api/submission/{problemName}/{id}", rp.apiReadSingleSubmission)

	go sp.Process()
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
