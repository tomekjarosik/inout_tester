package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

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

func (p *defaultSubmissionProcessor) processSubmission(submission SubmissionMetadata) (res SubmissionMetadata, err error) {
	fmt.Println("Processing submission:", submission)

	submission.State = Compiling
	p.store.Save(submission)

	compilationDir := path.Join(p.store.RootDir(), submission.ProblemName)
	solutionFilePath := path.Join(compilationDir, submission.SolutionFilename)
	executableFilePath := path.Join(compilationDir, submission.ExecutableFilename)
	submission.CompilationMode = testcase.AnalyzeGplusplusMode // TODO: pass from outside
	submission.CompilationOutput, err = testcase.CompileSolution(solutionFilePath, submission.CompilationMode, executableFilePath)

	if err != nil {
		submission.State = CompilationError
		p.store.Save(submission)
		return submission, err
	}
	defer os.Remove(executableFilePath)

	submission.State = RunningTests
	p.store.Save(submission)

	testcases, err := testcase.Populate(path.Join(p.problemsDataDirectory, submission.ProblemName))
	if err != nil {
		return submission, err
	}

	runner := testcase.NewRunner(p.problemsDataDirectory,
		testcase.DirectoryBasedDataStreamsProvider(path.Join(p.problemsDataDirectory, submission.ProblemName)))

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
