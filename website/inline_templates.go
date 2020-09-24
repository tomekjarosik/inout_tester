package website

import (
	"html/template"
	"time"
)

// TODO: maybe use proper templates??
// https://curtisvermeeren.github.io/2017/09/14/Golang-Templates-Cheatsheet#nesting
func HtmlHead() string {
	return `
<head>
	<!--Import Google Icon Font-->
	<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">

	<!-- Compiled and minified CSS -->
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
	
	<!--Let browser know website is optimized for mobile-->
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>

	<script>
	document.addEventListener('DOMContentLoaded', function() {
		var elems = document.querySelectorAll('.collapsible');
		var instances = M.Collapsible.init(elems,  {
			accordion: true
			});
	});
	document.addEventListener('DOMContentLoaded', function() {
		var elems = document.querySelectorAll('select');
		var instances = M.FormSelect.init(elems, {});
	  });
	</script>
</head>
	`
}

func HtmlDocumentWrap(content string) string {
	return `<!DOCTYPE html><html>` + content + `</html>`
}

func HomePageTemplate() (*template.Template, error) {
	return template.New("homepage").Funcs(
		template.FuncMap{
			"TimeFormat":          func(t time.Time) string { return t.Format(time.Stamp) },
			"ScoreColorFormat":    ScoreColorFormat,
			"TestCaseStatusColor": TestCaseStatusColorFormat,
		}).Parse(HtmlDocumentWrap(HtmlHead() + `
	<body class="container">
		<nav>
			<div class="nav-wrapper black">
			<a href="#" class="brand-logo">INOUT</a>
			<ul id="nav-mobile" class="right hide-on-med-and-down">
				<!-- <li><a href="/">Submissions</a></li>
				<li><a href="Ooooo">Components</a></li> -->
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
			<style type="text/css" scoped>
				td, th {
					border: 1px solid #dddddd;
					text-align: left;
					padding: 0px;
				}
			</style>
			<thead>
			<tr>
				<th>Test name</th>
				<th>Status</th>
				<th>Duration</th>
				<th>Additional info</th>
			</tr>
			<tbody>
			{{range .CompletedTestCases}}
				<tr class="{{TestCaseStatusColor .Result.Status}}">
					<td>{{.Info.Name}} </td>
					<td>{{.Result.Status}} </td>
					<td>{{.Result.Duration}} / {{.Info.TimeLimit}}</td>
					<td>{{.Result.Description}}</td>
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
`))
}

func SubmitForm() (*template.Template, error) {
	return template.New("submitForm").Parse(HtmlDocumentWrap(HtmlHead() + `
	<body class="container">

	<nav>
		<div class="nav-wrapper black">
		<a href="#" class="brand-logo">INOUT</a>
		<ul id="nav-mobile" class="right hide-on-med-and-down">
			<li><a href="/">Home</a></li>
		</ul>
		</div>
	</nav>

	<div class="divider"></div>

	<div class="section center-align">
	<form action="/api/submit" method="post" enctype="multipart/form-data">
		<div class="row input-field">
			<select name="problemName" required>
				<option value="" disabled selected>Choose problem</option>
				{{range .}}
				<option value="{{.Name}}">{{.Name}}</option>
				{{end}}
			</select>
	 	 </div>

		<div class="row input-field">
		  <select name="compilationMode" required>
			  <option value="" disabled selected>Choose compilation mode</option>
			  {{range .}}
			  <option value="{{.CompilationMode}}">{{.CompilationModeDescription}}</option>
			  {{end}}
		  </select>
		</div>

		<div class="row">
			<div class = "file-field input-field">
				<div class="btn">
					<span>Browse</span>
					<input type="file" name="solution" required/>
				</div>
				
				<div class="file-path-wrapper">
					<input class="file-path validate" type="text" name="solution" placeholder="Your solution .cpp"/>
				</div>
			</div>
		</div>
		<button class="btn waves-effect waves-light" type="submit" name="action">Submit
			<i class="material-icons right">send</i>
	  </button>
	</form>
	</div>
	
	<!--JavaScript at end of body for optimized loading-->
	<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
	</body>
`))
}
