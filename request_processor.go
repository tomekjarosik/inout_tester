package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/tomekjarosik/inout_tester/internal/submission"
	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
	"github.com/tomekjarosik/inout_tester/website"
)

// TODO(tjarosik): test g++/clang++ address sanitizers on Linux (+ write a README how to setup)
// TODO(tjarosik): add redirect to home page after submitting a solution
// TODO(tjarosik): add docker environment with adress sanitizers (g++ -lasan)

// RequestProcessor processes HTTP requests
type RequestProcessor struct {
	SubmissionStorage   submission.Storage
	SubmissionProcessor submission.Processor
}

// NewRequestProcessor constructor
func NewRequestProcessor(store submission.Storage, sp submission.Processor) RequestProcessor {
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
	compilationMode, _ := strconv.Atoi(r.Form.Get("compilationMode"))

	log.Println("compilationMode=", compilationMode)
	metadata := submission.NewMetadata(problemName, testcase.CompilationMode(compilationMode))
	fmt.Println("submissionMetadata:", metadata)
	err = rp.SubmissionProcessor.Submit(metadata, formFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, website.HtmlDocumentWrap(
		fmt.Sprintf(` <meta http-equiv="refresh" content="2;url=/" />
			File %s uploaded successfully! You will be redirected to the Home Page in 2 seconds...`, header.Filename)))
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

	problems, err := testcase.ListAvailableProblems("problems")
	if err != nil {
		http.Error(w, "failed read problems from 'problems' directory: "+err.Error(), http.StatusInternalServerError)
		return
	}
	compilationModes := []testcase.CompilationMode{testcase.ReleaseMode, testcase.AnalyzeClangMode, testcase.AnalyzeGplusplusMode}

	type ViewData struct {
		Problems         []string
		CompilationModes []testcase.CompilationMode
	}
	data := ViewData{Problems: problems, CompilationModes: compilationModes}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rp *RequestProcessor) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := website.HomePageTemplate()

	if err != nil {
		http.Error(w, "unable to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, rp.SubmissionStorage.List()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
