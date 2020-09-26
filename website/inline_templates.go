package website

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
