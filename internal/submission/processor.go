package submission

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// Processor processes submissions
type Processor interface {
	Submit(meta Metadata, sourceCode io.Reader) error
	Process() error
	Quit()
}

type defaultProcessor struct {
	queue           chan Metadata
	store           Storage
	testcaseArchive testcase.Archive
}

// NewProcessor constructor of the Processor
func NewProcessor(store Storage, testcaseArchive testcase.Archive) Processor {
	return &defaultProcessor{
		queue:           make(chan Metadata, 1000),
		store:           store,
		testcaseArchive: testcaseArchive,
	}
}

func (p *defaultProcessor) Submit(meta Metadata, solution io.Reader) error {
	fmt.Println("Submit:", meta)
	err := p.store.Upload(meta, solution)
	if err != nil {
		return err
	}
	p.queue <- meta

	return nil
}

func (p *defaultProcessor) processSubmission(submission Metadata) (res Metadata, err error) {
	fmt.Println("Processing submission:", submission)

	submission.Status = Compiling
	p.store.Save(submission)

	compilationDir := path.Join(p.store.RootDir(), submission.ProblemName)
	solutionFilePath := path.Join(compilationDir, submission.SolutionFilename)
	executableFilePath := path.Join(compilationDir, submission.ExecutableFilename)
	submission.CompilationOutput, err = testcase.CompileSolution(solutionFilePath, submission.CompilationMode, executableFilePath)

	if err != nil {
		submission.Status = CompilationError
		p.store.Save(submission)
		return submission, err
	}
	defer os.Remove(executableFilePath)

	submission.Status = RunningTests
	p.store.Save(submission)

	testcases, err := p.testcaseArchive.Testcases(submission.ProblemName)
	if err != nil {
		return submission, err
	}

	runner := p.testcaseArchive.Runner(submission.ProblemName)

	var processedTestCases []testcase.CompletedTestCase
	for _, tc := range testcases {
		executable := path.Join(compilationDir, submission.ExecutableFilename)
		res := runner.Run(executable, tc)
		processedTestCases = append(processedTestCases, testcase.CompletedTestCase{Info: tc, Result: res})
		submission.CompletedTestCases = processedTestCases
		p.store.Save(submission)
	}
	submission.CompletedTestCases = processedTestCases
	submission.Status = AllTestsCompleted
	err = p.store.Save(submission)
	log.Println("Processed submission", submission)
	return submission, err
}

func (p *defaultProcessor) Process() error {
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

func (p *defaultProcessor) Quit() {
	close(p.queue)
}
