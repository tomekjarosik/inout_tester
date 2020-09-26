package testcase

import (
	"io/ioutil"
	"strings"
	"time"
)

// PopulateTestCases searches directory for test case descriptions (.in / .out files, maybe others in the future)
// TODO: implement testcase metadata, which knows about timelimits and memory limits
func Populate(problemDataDir string) (testcases []Info, err error) {

	files, err := ioutil.ReadDir(problemDataDir)
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

func ListAvailableProblems(problemDataDir string) ([]string, error) {
	files, err := ioutil.ReadDir(problemDataDir)
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
