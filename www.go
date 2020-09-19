package main

import (
	"fmt"
	"net/http"
)

func wwwHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This will be home page")
}

// TODO(tjarosik): Read problem list from the Config
func wwwSubmitForm(w http.ResponseWriter, r *http.Request) {
	var form = `
	<!DOCTYPE html>
	<html>
	<body>
	
	<h1>The form element</h1>
	
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
