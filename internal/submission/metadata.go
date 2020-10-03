package submission

//go:generate stringer -type=Status

import (
	"encoding/json"
	"errors"
	"runtime"
	"time"

	guuid "github.com/google/uuid"
	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// ID unique identificator of a submission
type ID guuid.UUID

const metaFileExtension = ".meta"

// State submission status
type Status int

const (
	// Queued the solution is queued for processing
	Queued Status = iota + 1
	// Compiling the solution is being processed
	Compiling
	// CompilationError the solution failed to compile
	CompilationError
	// RunningTests currently running the provided test cases
	RunningTests
	// AllTestsCompleted all done
	AllTestsCompleted
)

// Metadata metadata of the submission
type Metadata struct {
	ID                  ID                           `json:"id"`
	SubmittedAt         time.Time                    `json:"submittedAt"`
	ProblemName         string                       `json:"problemName"`
	SolutionFilename    string                       `json:"solutionFilename"`
	Status              Status                       `json:"status"`
	ExecutableFilename  string                       `json:"executableFilename"`
	CompilationOutput   []byte                       `json:"compilationOutput"`
	CompilationMode     testcase.CompilationMode     `json:"compilationMode"`
	CompletedTestCases  []testcase.CompletedTestCase `json:"testCases"`
	TotalProcessingTime time.Duration                `json:"totalProcessingTime"`
	WorkerCount         int                          `json:"workerCount"`
}

func NewMetadata(problem string, mode testcase.CompilationMode) Metadata {
	id := ID(guuid.New())
	return Metadata{
		ID:                  id,
		SubmittedAt:         time.Now(),
		SolutionFilename:    id.String() + ".cpp",
		Status:              Queued,
		ProblemName:         problem,
		ExecutableFilename:  id.String() + ".tsk",
		CompilationMode:     mode,
		TotalProcessingTime: time.Duration(0),
		WorkerCount:         runtime.NumCPU() / 2,
	}
}

// TODO: Add tests for marshal / unmarshall
func (id ID) String() string {
	return guuid.UUID(id).String()
}

func (id ID) MarshalJSON() ([]byte, error) {
	g := guuid.UUID(id)
	return json.Marshal(g.String())
}
func (id *ID) UnmarshalJSON(data []byte) error {
	var g guuid.UUID
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
	(*id) = ID(g)
	return nil
}

func (e *Status) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	for i := 0; i <= len(_Status_index); i++ {
		if Status(i).String() == s {
			*e = Status(i)
			return nil
		}
	}
	return errors.New("invalid testacase status value")
}

func (e Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}
