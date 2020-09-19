package main

import (
	"fmt"
	"log"
	"net/http"
)

// RequestProcessor processes HTTP requests
type RequestProcessor struct {
	SubmissionProcessor SubmissionProcessor
}

// NewRequestProcessor constructor
func NewRequestProcessor(sp SubmissionProcessor) RequestProcessor {
	return RequestProcessor{sp}
}

func (rp *RequestProcessor) apiSubmitSolutionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Fatal(err)
	}

	metadata := NewSubmissionMetadata(r.Form.Get("problemName"))

	formFile, header, err := r.FormFile("solution")
	if err != nil {
		http.Error(w, "unable to open solution file", http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	err = rp.SubmissionProcessor.Submit(metadata, formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(w, "File Uploaded Successfully! ")
	fmt.Fprintln(w, "Name of the File: ", header.Filename)
	fmt.Fprintln(w, "Size of the File: ", header.Size)
}

func (rp *RequestProcessor) apiReadSingleSubmission(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]
	//problemName := vars["problemName"]
}
