package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tomekjarosik/inout_tester/website"
)

// TODO(tjarosik): pass compilation mode properly
// TODO(tjarosik)

// RequestProcessor processes HTTP requests
type RequestProcessor struct {
	SubmissionStorage   SubmissionStorage
	SubmissionProcessor SubmissionProcessor
}

// NewRequestProcessor constructor
func NewRequestProcessor(store SubmissionStorage, sp SubmissionProcessor) RequestProcessor {
	return RequestProcessor{store, sp}
}

func (rp *RequestProcessor) apiSubmitSolutionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Fatal(err)
	}

	formFile, header, err := r.FormFile("solution")
	if err != nil {
		http.Error(w, "unable to open solution file", http.StatusBadRequest)
		return
	}
	defer formFile.Close()
	problemName := r.Form.Get("problemName")
	compilationMode := r.Form.Get("compilationMode")
	log.Println("compilationMode=", compilationMode)
	err = rp.SubmissionProcessor.Submit(problemName, formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

// TODO(tjarosik): Read problem list from the Config
func (rp *RequestProcessor) wwwSubmitForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := website.SubmitForm()

	if err != nil {
		http.Error(w, "unable to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	type Data struct {
		Name                       string
		CompilationMode            string
		CompilationModeDescription string
	}
	problems := []Data{
		{"volvo", "cpp_release", "c++ release"},
		{"saab", "cpp_analyze", "c++ analyze"}}
	if err = tmpl.Execute(w, problems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rp *RequestProcessor) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := website.HomePageTemplate()

	if err != nil {
		http.Error(w, "unable to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	submissions, err := rp.SubmissionStorage.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, submissions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
