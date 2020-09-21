package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

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

	err = rp.SubmissionProcessor.Submit(r.Form.Get("problemName"), formFile)
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
	var form = `
	<!DOCTYPE html>
	<html>
	<body>
	
	<h1>Submit your solution</h1>
	
	<form action="/api/submit" method="post" enctype="multipart/form-data">
	<label for="problemName">Choose a problem:</label>
		<select name="problemName" id="problemName">
		<option value="volvo">volvo</option>
		<option value="saab">saab</option>
		</select>
		<br><br>
	  <label for="myfile">Select a file:</label>
		<input type="file" id="solution" name="solution">
		<input type="submit" value="Submit">
	</form>
	
	<p>Click the "Submit" button and the form-data will be sent".</p>
	
	</body>
	</html>
`
	fmt.Fprintln(w, form)
}

// TODO: Ekrany
// Main -> Moje zgłoszenia:
//         "Czas zgłoszenia"	"Zadanie"  "Status"	Wynik (liczba testow OK/All)

func (rp *RequestProcessor) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("homepage").Parse(`
	<style type="text/css">
.tg  {border-collapse:collapse;border-spacing:0;margin:0px auto;}
.tg td{border-color:black;border-style:solid;border-width:1px;font-family:Arial, sans-serif;font-size:14px;
  overflow:hidden;padding:10px 5px;word-break:normal;}
.tg th{border-color:black;border-style:solid;border-width:1px;font-family:Arial, sans-serif;font-size:14px;
  font-weight:normal;overflow:hidden;padding:10px 5px;word-break:normal;}
.tg .tg-1wig{font-weight:bold;text-align:left;vertical-align:top}
.tg .tg-0lax{text-align:left;vertical-align:top}
@media screen and (max-width: 767px) {.tg {width: auto !important;}.tg col {width: auto !important;}.tg-wrap {overflow-x: auto;-webkit-overflow-scrolling: touch;margin: auto 0px;}}</style>
<div class="tg-wrap"><table class="tg">
<thead>
  <tr>
    <th class="tg-1wig">Submitted at</th>
    <th class="tg-1wig">Problem</th>
    <th class="tg-1wig">Status</th>
    <th class="tg-1wig">Score</th>
  </tr>
</thead>
<tbody>
{{range .}}
  <tr>
    <td class="tg-0lax">{{.SubmittedAt}} </td>
    <td class="tg-0lax">{{.ProblemName}} </td>
    <td class="tg-0lax">{{.State}}</td>
    <td class="tg-0lax">{{.Score }}</td>
  </tr>
{{end}}
</tbody>
</table></div>
`)

	if err != nil {
		http.Error(w, "unable to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	submissions, err := rp.SubmissionStorage.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, submissions)
}
