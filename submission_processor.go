package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"

	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// SubmissionProcessor processes submissions
type SubmissionProcessor interface {
	Submit(problemName string, sourceCode io.Reader) error
	Process() error
	Quit()
}

type defaultSubmissionProcessor struct {
	queue                 chan SubmissionMetadata
	store                 SubmissionStorage
	problemsDataDirectory string
}

// NewSubmissionProcessor constructor of SubmissionProcessor
func NewSubmissionProcessor(store SubmissionStorage) SubmissionProcessor {
	return &defaultSubmissionProcessor{
		queue:                 make(chan SubmissionMetadata, 1000),
		store:                 store,
		problemsDataDirectory: "problems",
	}
}

func (p *defaultSubmissionProcessor) Submit(problemName string, solution io.Reader) error {
	log.Println("defaultSubmissionProcessor Submit()")
	metadata, err := p.store.Upload(problemName, solution)
	if err != nil {
		return err
	}
	p.queue <- metadata

	return nil
}

// TODO: Implement (TestCase)Provider in testcase package, which can populate test cases
func (p *defaultSubmissionProcessor) PopulateTestCases(problemDataDir string) (testcases []testcase.Info, err error) {

	files, err := ioutil.ReadDir(problemDataDir)
	if err != nil {
		return
	}
	const ext = ".in"
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ext) {
			continue
		}
		tc := testcase.NewInfo(strings.TrimSuffix(f.Name(), ext), 10*time.Second, 0)
		testcases = append(testcases, tc)
	}
	return
}

func (p *defaultSubmissionProcessor) processSubmission(submission SubmissionMetadata) (SubmissionMetadata, error) {
	fmt.Println("Processing submission:", submission)
	testcases, err := p.PopulateTestCases(path.Join(p.problemsDataDirectory, submission.ProblemName))
	if err != nil {
		return submission, err
	}
	submission.State = Compiling
	p.store.Save(submission)

	compilationDir := path.Join(p.store.RootDir(), submission.ProblemName)
	submission.CompilationOutput, err = testcase.CompileSolution(
		path.Join(compilationDir, submission.SolutionFilename),
		path.Join(compilationDir, submission.ExecutableFilename))

	if err != nil {
		submission.State = CompilationError
		p.store.Save(submission)
		return submission, err
	}
	submission.State = RunningTests
	p.store.Save(submission)

	// TODO: name?
	runner := testcase.NewRunner(p.problemsDataDirectory, path.Join(p.problemsDataDirectory, submission.ProblemName))

	var processedTestCases []testcase.CompletedTestCase
	for _, tc := range testcases {
		executable := path.Join(compilationDir, submission.ExecutableFilename)
		res := runner.Run(executable, tc)
		processedTestCases = append(processedTestCases, testcase.CompletedTestCase{Info: tc, Result: res})
	}
	submission.CompletedTestCases = processedTestCases
	submission.State = RunAllTests
	err = p.store.Save(submission)
	log.Println("Processed submission", submission)
	return submission, err
}

func (p *defaultSubmissionProcessor) Process() error {
	if err := p.store.LoadAll(); err != nil {
		log.Panic(err)
	}
	for submission := range p.queue {
		_, err := p.processSubmission(submission)
		if err != nil {
			log.Println("ProcessSubmission returned error: ", err)
		}
	}
	fmt.Println("defaultSubmissionProcessor has exited successfully.")
	return nil
}

func (p *defaultSubmissionProcessor) Quit() {
	close(p.queue)
}
