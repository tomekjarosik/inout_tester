package website

import (
	"html/template"

	"github.com/tomekjarosik/inout_tester/internal/testcase"
)

func SubmitForm() (*template.Template, error) {
	return template.New("submitForm").Funcs(template.FuncMap{
		"AsInt":                    func(mode testcase.CompilationMode) int { return int(mode) },
		"FullCompilationCommadFor": testcase.FullCompilationCommadFor,
	}).Parse(HtmlDocumentWrap(HtmlHead() + `
	<body class="container">

	<nav>
		<div class="nav-wrapper black">
		<a href="/" class="brand-logo">INOUT</a>
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
				{{range .Problems}}
				<option value="{{.}}">{{.}}</option>
				{{end}}
			</select>
	 	 </div>

		<div class="row input-field">
		  <select name="compilationMode" required>
			  <option value="" disabled selected>Choose compilation mode</option>
			  {{range .CompilationModes}}
			  <option value="{{AsInt .}}">{{.}} [{{FullCompilationCommadFor .}}]</span></option>
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
