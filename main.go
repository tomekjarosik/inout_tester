package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/tomekjarosik/inout_tester/internal/submission"
	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

var flagPort int
var flagProblemsDirectory string
var flagSubmissionsDirectory string

func init() {
	flag.IntVar(&flagPort, "port", 8080, "Webserver port")
	flag.StringVar(&flagProblemsDirectory, "problems-dir", "problems",
		"Root directory where problems are located. Each problem is a sub-dir and contains test data (.in/.out files)")
	flag.StringVar(&flagSubmissionsDirectory, "submissions-dir", "submissions", "Directory where submissions will be stored")
}

func generateMultiplyBy2(dir string) {
	assert := func(e error) {
		if e != nil {
			log.Panic(e)
		}
	}
	problemPath := path.Join(dir, "multiply_by_2")

	assert(os.MkdirAll(problemPath, 0755))

	assert(ioutil.WriteFile(path.Join(problemPath, "t1.in"), []byte("1\n"), 0666))
	assert(ioutil.WriteFile(path.Join(problemPath, "t1.out"), []byte("2\n"), 0666))

	assert(ioutil.WriteFile(path.Join(problemPath, "t2.in"), []byte("123\n"), 0666))
	assert(ioutil.WriteFile(path.Join(problemPath, "t2.out"), []byte("246\n"), 0666))

	assert(ioutil.WriteFile(path.Join(problemPath, "t3.in"), []byte("-2000000000\n"), 0666))
	assert(ioutil.WriteFile(path.Join(problemPath, "t3.out"), []byte("-4000000000\n"), 0666))

	assert(ioutil.WriteFile(path.Join(problemPath, "t4.in"), []byte("-99999999999999999999999999999999999999999999999999\n"), 0666))
	assert(ioutil.WriteFile(path.Join(problemPath, "t4.out"), []byte("-199999999999999999999999999999999999999999999999998\n"), 0666))
}

// TODO: Handle Ctrl+C properly
// TODO: add ability to run tests in parallel, for each submission
func main() {
	fmt.Println("Starting...")
	flag.Parse()

	storage := submission.NewDefaultStorage(flagSubmissionsDirectory)
	if err := storage.Init(); err != nil {
		log.Panic(err)
	}
	testcaseArchive := testcase.NewArchive(flagProblemsDirectory)
	sp := submission.NewProcessor(storage, testcaseArchive)
	rp := NewRequestProcessor(storage, sp, testcaseArchive)

	problems, err := testcaseArchive.Problems()
	if err != nil || len(problems) == 0 {
		generateMultiplyBy2(flagProblemsDirectory)
	}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", rp.RenderHomePage)
	myRouter.HandleFunc("/submit", rp.wwwSubmitForm)
	myRouter.HandleFunc("/api/submit", rp.apiSubmitSolutionHandler).Methods("POST")
	myRouter.HandleFunc("/api/submission/{problemName}/{id}", rp.apiReadSingleSubmission)

	go sp.Process()
	fmt.Printf("Started new server at http://localhost:%d\n", flagPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flagPort), myRouter))
}
