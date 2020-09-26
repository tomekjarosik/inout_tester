package website

import (
	"html/template"
	"time"

	"github.com/tomekjarosik/inout_tester/internal/testcase"
)

func HomePageTemplate() (*template.Template, error) {
	return template.New("homepage").Funcs(
		template.FuncMap{
			"TimeFormat":          func(t time.Time) string { return t.Format(time.Stamp) },
			"ScoreColorFormat":    ScoreColorFormat,
			"TestCaseStatusColor": TestCaseStatusColorFormat,
			"HasAnyTestCases":     func(c []testcase.CompletedTestCase) bool { return len(c) > 0 },
			"BytesToString":       func(arr []byte) string { return string(arr) },
			"FullCompilationCommandFor": func(cm testcase.CompilationMode) string {
				cmd, err := testcase.CompilationCommand("a.cpp", cm, "a.out")
				if err != nil {
					return "unable to convert"
				}
				return cmd.String()
			},
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
		<div class="collapsible-header">
			<span style="font-weight:bold">{{TimeFormat .SubmittedAt}}&nbsp;|&nbsp;</span>{{.ProblemName}}</span>&nbsp;&nbsp;{{.Status}} 
			<span class="new badge {{ScoreColorFormat .Score}}" data-badge-caption="points">{{.Score}}/{{.MaxScore}}</span>
		</div>
		<div class="collapsible-body">
			<div style="border: 2px solid black; background: lightblue;">
			{{FullCompilationCommandFor .CompilationMode}}
			</div>
			{{if HasAnyTestCases .CompletedTestCases}}
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
			{{else}}
			<div style="border: 2px solid red;">
			<p>{{BytesToString .CompilationOutput}}</p>
			</div>
			{{end}}
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
