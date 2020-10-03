package submission

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

type inMemoryRunner struct {
}

type imMemoryArchive struct {
}

func NewInMemoryArchive() testcase.Archive {
	return &imMemoryArchive{}
}

func (runner *inMemoryRunner) Run(executable string, info testcase.Info) testcase.Result {
	return testcase.Result{Status: testcase.Accepted, Description: info.Name}
}

func (archive *imMemoryArchive) Problems() ([]string, error) {
	return []string{"problem1"}, nil
}
func (archive *imMemoryArchive) Testcases(problemName string) (testcases []testcase.Info, err error) {
	return []testcase.Info{
		testcase.Info{Name: "t20"},
		testcase.Info{Name: "t10"},
		testcase.Info{Name: "t15"},
		testcase.Info{Name: "t13"},
		testcase.Info{Name: "t17"},
	}, nil
}

func (archive *imMemoryArchive) Runner(problemName string) testcase.Runner {
	return &inMemoryRunner{}
}

func TestProcessor_ProcessSolution(t *testing.T) {

	dirname, err := ioutil.TempDir(os.TempDir(), "testprocessor-*")
	assert.NoError(t, err)
	storage := NewDefaultStorage(dirname)
	storage.Init()

	testcaseArchive := NewInMemoryArchive()
	proc := NewProcessor(storage, testcaseArchive)

	metadata := NewMetadata("problem1", testcase.ReleaseMode)
	sol := strings.NewReader(`#include <cstdio>
	int main() { printf("1\n"); return 0; }
	`)
	storage.Upload(metadata, sol)

	assert.Equal(t, Queued, metadata.Status)

	go proc.Process()
	proc.Submit(metadata)
	// TODO: wait smarter
	time.Sleep(750 * time.Millisecond)

	metadata, ok := storage.Get(metadata.ID)
	assert.True(t, ok)
	assert.Equal(t, AllTestsCompleted, metadata.Status)
	assert.Equal(t, "t10", metadata.CompletedTestCases[0].Info.Name)
	assert.Equal(t, "t10", metadata.CompletedTestCases[0].Result.Description)
	assert.Equal(t, testcase.Accepted, metadata.CompletedTestCases[0].Result.Status)

	assert.Equal(t, "t13", metadata.CompletedTestCases[1].Info.Name)
	assert.Equal(t, "t15", metadata.CompletedTestCases[2].Info.Name)
	assert.Equal(t, "t17", metadata.CompletedTestCases[3].Info.Name)
	assert.Equal(t, "t20", metadata.CompletedTestCases[4].Info.Name)

	proc.Quit()
	storage.Destroy()
	// for extra safety
	if strings.Contains(dirname, "testprocessor-") {
		os.RemoveAll(dirname)
	}
}
