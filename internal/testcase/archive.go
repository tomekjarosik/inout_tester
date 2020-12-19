package testcase

import (
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Archive interface {
	Problems() ([]string, error)
	Testcases(problemName string) (testcases []Info, err error)
	Runner(problemName string) Runner
}

type defaultArchive struct {
	dataDir string
}

func NewArchive(problemsDirectory string) Archive {
	return &defaultArchive{
		dataDir: problemsDirectory,
	}
}

// ByTestcaseName is a helper type to implement sorting
type ByTestcaseStatusAndName []CompletedTestCase

func (a ByTestcaseStatusAndName) Len() int      { return len(a) }
func (a ByTestcaseStatusAndName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTestcaseStatusAndName) Less(i, j int) bool {
	if a[i].Result.Status == a[j].Result.Status {
		return a[i].Info.Name < a[j].Info.Name
	}
	if a[j].Result.Status == Accepted {
		return true
	}
	if a[i].Result.Status == Accepted {
		return false
	}
	return a[i].Info.Name < a[j].Info.Name
}

// Testcases searches directory for test case descriptions (.in / .out files, maybe others in the future)
// TODO: implement testcase metadata, which knows about timelimits and memory limits
func (a *defaultArchive) Testcases(problemName string) (testcases []Info, err error) {

	files, err := ioutil.ReadDir(path.Join(a.dataDir, problemName))
	if err != nil {
		return
	}
	const ext = ".in"
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ext) {
			continue
		}
		tc := NewInfo(strings.TrimSuffix(f.Name(), ext), 10*time.Second, 0)
		testcases = append(testcases, tc)
	}
	return
}

func (a *defaultArchive) Problems() ([]string, error) {
	files, err := ioutil.ReadDir(a.dataDir)
	if err != nil {
		return []string{""}, err
	}
	res := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			res = append(res, f.Name())
		}
	}
	return res, nil
}

func (a *defaultArchive) Runner(problemName string) Runner {
	return NewRunner(problemName,
		DirectoryBasedDataStreamsProvider(path.Join(a.dataDir, problemName)))
}
