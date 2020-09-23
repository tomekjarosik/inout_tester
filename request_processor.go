package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// TODO: Use Material https://materializecss.com/navbar.html for all pages
// TODO: Refactor sa MaterialCSS can be easily used e.g extract header and common stuff
//

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

// TODO: Special function to colorize row by "Accepted/Wrong Answer/TimeLimitExceeded"
func (rp *RequestProcessor) RenderHomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("homepage").Funcs(
		template.FuncMap{
			"TimeFormat": func(t time.Time) string { return t.Format(time.Stamp) },
			"ScoreColorFormat": func(score int) string {
				if score == 0 {
					return "red"
				} else {
					return "blue"
				}
			},
			"TestCasesStatusColor": func(status testcase.Status) string {
				if status == testcase.Accepted {
					return "green lighten-3"
				} else {
					return " red lighten-3"
				}
			},
		}).Parse(`
<!DOCTYPE html>
<html>
	<head>
		<!--Import Google Icon Font-->
		<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

		<!-- Compiled and minified CSS -->
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
		
		<!-- Compiled and minified JavaScript -->
		
		<!--Let browser know website is optimized for mobile-->
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>

		<script>
		document.addEventListener('DOMContentLoaded', function() {
			var elems = document.querySelectorAll('.collapsible');
			var instances = M.Collapsible.init(elems,  {
				accordion: true
			  });
		});
		</script>
		<style>
		td, th {
			border: 1px solid #dddddd;
			text-align: left;
			padding: 0px;
		  }
		</style>
	</head>

	<body class = "container">
		<nav>
			<div class="nav-wrapper black">
			<a href="#" class="brand-logo">INOUT</a>
			<ul id="nav-mobile" class="right hide-on-med-and-down">
				<li><a href="/">Submissions</a></li>
				<li><a href="Ooooo">Components</a></li>
				<li><a class="waves-effect waves-light btn" href="/submit"><i class="material-icons right">cloud_upload</i>Submit</a></li>
			</ul>
			</div>
		</nav>

		<div class="divider"></div>
		
		<div class="section center-align">
		{{range .}}
		<ul class="collapsible">
		<li>
		<div class="collapsible-header"> {{TimeFormat .SubmittedAt}} for {{.ProblemName}} {{.State}} <span class="new badge {{ScoreColorFormat .Score}}" data-badge-caption="points">{{.Score}}</span> </div>
		<div class="collapsible-body">
			<table class="responsive-table striped" cellspacing="0">
			<thead>
			<tr>
				<th>Test name</th>
				<th>Status</th>
				<th>Duration</th>
				<th>Additional info</th>
			</tr>
			<tbody>
			{{range .CompletedTestCases}}
				<tr class="{{TestCasesStatusColor .Result.Status}}">
					<td>{{.Info.Name}} </td>
					<td>{{.Result.Status}} </td>
					<td>{{.Result.Duration}} / {{.Info.TimeLimit}}</td>
					<td>{{.Result.StatusDescription}}</td>
				</tr>
			{{end}}
			</tbody> 
			</table>
			</div>
		</li>
		</ul>
		{{end}}
		</div>
	<!--JavaScript at end of body for optimized loading-->
	<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
	</body>
</html>
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
