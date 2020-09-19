package main

import (
	"path"
	"time"

	guuid "github.com/google/uuid"
)

// SubmissionState submission state
type SubmissionState int

const (
	// Queued the solution is queued for processing
	Queued SubmissionState = iota
	// Compiling the solution is being processed
	Compiling
	// Processing the solution is being processed
	Processing
	// Processed the solution has been processed and results are available
	Processed
)

// TesCaseStatus status of the
type TesCaseStatus int

const (
	// NotRunYet test is waiting to be run
	NotRunYet TesCaseStatus = iota
	// InternalError something unexpected went wrong
	InternalError
	// CompilationError the solution failed to compile
	CompilationError
	// TimeLimitExceeded the test took too long to process
	TimeLimitExceeded
	// MemoryLimitExceeded the test run used too much RAM
	MemoryLimitExceeded
	// WrongAnswer test run successfully but test outputs differ
	WrongAnswer
	// TestCaseOK all OK
	TestCaseOK
)

// TestCase struct describing results of a single test run (input, output, expected output)
type TestCase struct {
	Name                 string        `json:"name"`
	InputFilename        string        `json:"in"`
	GoldenOutputFilename string        `json:"goldenOut"`
	OutputFilename       string        `json:"out"`
	TimeLimit            time.Duration `json:"timeLimit"`
	MemoryLimit          int           `json:"memoryLimit"`

	Status            TesCaseStatus `json:"status"`
	StatusDescription string        `json:"statusDescription"`
	Duration          time.Duration `json:"duration"`
}

// NewTestCase construct of TestCase struct
func NewTestCase(solutionID string, testName string, problemDataDir string, outputDir string) TestCase {
	return TestCase{
		Name:                 testName,
		Status:               NotRunYet,
		Duration:             0,
		InputFilename:        path.Join(problemDataDir, testName+".in"),
		GoldenOutputFilename: path.Join(problemDataDir, testName+".out"),
		OutputFilename:       path.Join(outputDir, testName+"."+solutionID+".out"),
		TimeLimit:            10 * time.Second,
	}
}

// SubmissionMetadata metadata of the submission
type SubmissionMetadata struct {
	ID                 string          `json:"id"`
	SubmittedAt        time.Time       `JSON:"submittedAt"`
	ProblemName        string          `json:"problemName"`
	SolutionFilename   string          `json:"solutionFilename"`
	State              SubmissionState `json:"state"`
	ExecutableFilename string          `json:"executableFilename"`
	CompilationOutput  []byte          `json:"compilationOutput"`
	TestCases          []TestCase      `json:"testCases"`
}

// NewSubmissionMetadata create new SubmissionMetadata object with unique ID
func NewSubmissionMetadata(problemName string) SubmissionMetadata {
	id := guuid.New().String()
	return SubmissionMetadata{
		ID:                 id,
		SubmittedAt:        time.Now(),
		SolutionFilename:   id + ".cpp",
		State:              Queued,
		ProblemName:        problemName,
		ExecutableFilename: id + ".tsk",
	}
}
