package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// SubmissionProcessor processes submissions
type SubmissionProcessor interface {
	Submit(metadata SubmissionMetadata, sourceCode io.Reader) error
	Process() error
	Quit()
}

type defaultSubmissionProcessor struct {
	queue                 chan SubmissionMetadata
	quit                  chan int
	submissionsDirectory  string
	queueDirectory        string
	problemsDataDirectory string
}

// NewSubmissionProcessor constructor of SubmissionProcessor
func NewSubmissionProcessor() SubmissionProcessor {
	return &defaultSubmissionProcessor{
		queue:                 make(chan SubmissionMetadata),
		submissionsDirectory:  "submissions",
		queueDirectory:        "queue",
		problemsDataDirectory: "problems",
	}
}

func ensureDirectoryExists(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func (p *defaultSubmissionProcessor) saveSubmission(metadata SubmissionMetadata, directory string) error {
	metadataFile, err := os.Create(path.Join(directory, metadata.ID))
	if err != nil {
		return err
	}
	defer metadataFile.Close()
	enc := json.NewEncoder(metadataFile)
	enc.SetIndent("", "\t")
	return enc.Encode(metadata)
}

func (p *defaultSubmissionProcessor) deleteSubmission(metadata SubmissionMetadata, directory string) error {
	return os.Remove(path.Join(directory, metadata.ID))
}

func (p *defaultSubmissionProcessor) Submit(metadata SubmissionMetadata, solution io.Reader) error {
	log.Println("defaultSubmissionProcessor Submit()")
	solutionsDirectory := path.Join(p.submissionsDirectory, metadata.ProblemName)

	submissionFile, err := os.Create(path.Join(solutionsDirectory, metadata.SolutionFilename))
	if err != nil {
		return err
	}
	defer submissionFile.Close()

	_, err = io.Copy(submissionFile, solution)
	if err != nil {
		return err
	}
	if err = p.saveSubmission(metadata, p.queueDirectory); err != nil {
		return err
	}
	p.queue <- metadata

	return nil
}

func (p *defaultSubmissionProcessor) ensureDirectoryStructure() error {
	if err := ensureDirectoryExists(p.submissionsDirectory); err != nil {
		return err
	}
	if err := ensureDirectoryExists(p.queueDirectory); err != nil {
		return err
	}
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
		path.Join(p.submissionsDirectory, submission.ProblemName), submission.ID)
	if err != nil {
		return submission, err
	}
	submission.State = Compiling
	compilationDir := path.Join(p.submissionsDirectory, submission.ProblemName)
	submission.CompilationOutput, err = compileSolution(
		path.Join(compilationDir, submission.SolutionFilename),
		path.Join(compilationDir, submission.ExecutableFilename))
	if err != nil {
		return submission, err
	}
	var processedTestCases []TestCase
	for _, tc := range testcases {
		processedTc := runSingleTestCase(path.Join(compilationDir, submission.ExecutableFilename), tc)
		processedTestCases = append(processedTestCases, processedTc)
	}
	submission.TestCases = processedTestCases
	err = p.saveSubmission(submission, path.Join(p.submissionsDirectory, submission.ProblemName))
	log.Println("Processed submission", submission)
	return submission, err
}

func (p *defaultSubmissionProcessor) Process() error {
	if err := p.ensureDirectoryStructure(); err != nil {
		return err
	}
	for submission := range p.queue {
		p.deleteSubmission(submission, p.queueDirectory)
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
