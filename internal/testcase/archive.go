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
