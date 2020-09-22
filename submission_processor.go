package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
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

func (p *defaultSubmissionProcessor) PopulateTestCases(problemDataDir string, outputDir string, solutionID string) (testcases []TestCase, err error) {

	files, err := ioutil.ReadDir(problemDataDir)
	if err != nil {
		return
	}
	const ext = ".in"
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ext) {
			continue
		}
		tc := NewTestCase(solutionID, strings.TrimSuffix(f.Name(), ext), problemDataDir, outputDir)
		testcases = append(testcases, tc)
	}
	return
}

func (p *defaultSubmissionProcessor) processSubmission(submission SubmissionMetadata) (SubmissionMetadata, error) {
	fmt.Println("Processing submission:", submission)
	testcases, err := p.PopulateTestCases(path.Join(p.problemsDataDirectory, submission.ProblemName),
		path.Join(p.store.RootDir(), submission.ProblemName), submission.ID.String())
	if err != nil {
		return submission, err
	}
	submission.State = Compiling
	p.store.Save(submission)

	compilationDir := path.Join(p.store.RootDir(), submission.ProblemName)
	submission.CompilationOutput, err = compileSolution(
		path.Join(compilationDir, submission.SolutionFilename),
		path.Join(compilationDir, submission.ExecutableFilename))

	if err != nil {
		submission.State = CompilationError
		p.store.Save(submission)
		return submission, err
	}
	submission.State = RunningTests
	p.store.Save(submission)

	var processedTestCases []TestCase
	for _, tc := range testcases {
		processedTc := runSingleTestCase(path.Join(compilationDir, submission.ExecutableFilename), tc)
		processedTestCases = append(processedTestCases, processedTc)
	}
	submission.TestCases = processedTestCases
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
